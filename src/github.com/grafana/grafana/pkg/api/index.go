package api

import (
	"fmt"
	"strings"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/setting"
)

func setIndexViewData(c *m.ReqContext) (*dtos.IndexViewData, error) {
	settings, err := getFrontendSettingsMap(c)
	if err != nil {
		return nil, err
	}

	prefsQuery := m.GetPreferencesWithDefaultsQuery{OrgId: c.OrgId, UserId: c.UserId}
	if err := bus.Dispatch(&prefsQuery); err != nil {
		return nil, err
	}
	prefs := prefsQuery.Result

	// Read locale from acccept-language
	acceptLang := c.Req.Header.Get("Accept-Language")
	locale := "ko-KR"

	if len(acceptLang) > 0 {
		parts := strings.Split(acceptLang, ",")
		locale = parts[0]
	}

	appURL := setting.AppUrl
	appSubURL := setting.AppSubUrl

	// special case when doing localhost call from phantomjs
	if c.IsRenderCall {
		appURL = fmt.Sprintf("%s://localhost:%s", setting.Protocol, setting.HttpPort)
		appSubURL = ""
		settings["appSubUrl"] = ""
	}

	hasEditPermissionInFoldersQuery := m.HasEditPermissionInFoldersQuery{SignedInUser: c.SignedInUser}
	if err := bus.Dispatch(&hasEditPermissionInFoldersQuery); err != nil {
		return nil, err
	}

	var data = dtos.IndexViewData{
		User: &dtos.CurrentUser{
			Id:                         c.UserId,
			IsSignedIn:                 c.IsSignedIn,
			Login:                      c.Login,
			Email:                      c.Email,
			Name:                       c.Name,
			OrgCount:                   c.OrgCount,
			OrgId:                      c.OrgId,
			OrgName:                    c.OrgName,
			OrgRole:                    c.OrgRole,
			GravatarUrl:                dtos.GetGravatarUrl(c.Email),
			IsGrafanaAdmin:             c.IsGrafanaAdmin,
			LightTheme:                 prefs.Theme == "light",
			Timezone:                   prefs.Timezone,
			Locale:                     locale,
			HelpFlags1:                 c.HelpFlags1,
			HasEditPermissionInFolders: hasEditPermissionInFoldersQuery.Result,
		},
		Settings:                settings,
		Theme:                   prefs.Theme,
		AppUrl:                  appURL,
		AppSubUrl:               appSubURL,
		GoogleAnalyticsId:       setting.GoogleAnalyticsId,
		GoogleTagManagerId:      setting.GoogleTagManagerId,
		BuildVersion:            setting.BuildVersion,
		BuildCommit:             setting.BuildCommit,
		NewGrafanaVersion:       plugins.GrafanaLatestVersion,
		NewGrafanaVersionExists: plugins.GrafanaHasUpdate,
		AppName:                 setting.ApplicationName,
	}

	if setting.DisableGravatar {
		data.User.GravatarUrl = setting.AppSubUrl + "/public/img/user_profile.png"
	}

	if len(data.User.Name) == 0 {
		data.User.Name = data.User.Login
	}

	themeURLParam := c.Query("theme")
	if themeURLParam == "light" {
		data.User.LightTheme = true
		data.Theme = "light"
	}

	if hasEditPermissionInFoldersQuery.Result {
		children := []*dtos.NavLink{
			{Text: "대쉬보드", Icon: "gicon gicon-dashboard-new", Url: setting.AppSubUrl + "/dashboard/new"},
		}

		if c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR {
			children = append(children, &dtos.NavLink{Text: "폴더", SubTitle: "대쉬보드 운영을 위한 폴더 생성", Id: "folder", Icon: "gicon gicon-folder-new", Url: setting.AppSubUrl + "/dashboards/folder/new"})
		}

		children = append(children, &dtos.NavLink{Text: "임포트", SubTitle: "대쉬보드 임포트", Id: "import", Icon: "gicon gicon-dashboard-import", Url: setting.AppSubUrl + "/dashboard/import"})

		data.NavTree = append(data.NavTree, &dtos.NavLink{
			Text:     "생성",
			Id:       "create",
			Icon:     "fa fa-fw fa-plus",
			Url:      setting.AppSubUrl + "/dashboard/new",
			Children: children,
		})
	}

	dashboardChildNavs := []*dtos.NavLink{
		{Text: "홈", Id: "home", Url: setting.AppSubUrl + "/", Icon: "gicon gicon-home", HideFromTabs: true},
		{Text: "분리", Divider: true, Id: "divider", HideFromTabs: true},
		{Text: "관리", Id: "manage-dashboards", Url: setting.AppSubUrl + "/dashboards", Icon: "gicon gicon-manage"},
		{Text: "재생목록", Id: "playlists", Url: setting.AppSubUrl + "/playlists", Icon: "gicon gicon-playlists"},
		{Text: "스냅샷", Id: "snapshots", Url: setting.AppSubUrl + "/dashboard/snapshots", Icon: "gicon gicon-snapshots"},
	}

	data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:     "대쉬보드",
		Id:       "dashboards",
		SubTitle: "폴더 및 대쉬보드 관리",
		Icon:     "gicon gicon-dashboard",
		Url:      setting.AppSubUrl + "/",
		Children: dashboardChildNavs,
	})

	if setting.ExploreEnabled && (c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR) {
		data.NavTree = append(data.NavTree, &dtos.NavLink{
			Text:     "탐색",
			Id:       "explore",
			SubTitle: "데이터 탐색",
			Icon:     "fa fa-rocket",
			Url:      setting.AppSubUrl + "/explore",
			Children: []*dtos.NavLink{
				{Text: "New tab", Icon: "gicon gicon-dashboard-new", Url: setting.AppSubUrl + "/explore"},
			},
		})
	}

	if c.IsSignedIn {
		// Only set login if it's different from the name
		var login string
		if c.SignedInUser.Login != c.SignedInUser.NameOrFallback() {
			login = c.SignedInUser.Login
		}
		profileNode := &dtos.NavLink{
			Text:         c.SignedInUser.NameOrFallback(),
			SubTitle:     login,
			Id:           "profile",
			Img:          data.User.GravatarUrl,
			Url:          setting.AppSubUrl + "/profile",
			HideFromMenu: true,
			Children: []*dtos.NavLink{
				{Text: "설정", Id: "profile-settings", Url: setting.AppSubUrl + "/profile", Icon: "gicon gicon-preferences"},
				{Text: "비밀번호 변경", Id: "change-password", Url: setting.AppSubUrl + "/profile/password", Icon: "fa fa-fw fa-lock", HideFromMenu: true},
			},
		}

		if !setting.DisableSignoutMenu {
			// add sign out first
			profileNode.Children = append(profileNode.Children, &dtos.NavLink{
				Text: "로그아웃", Id: "sign-out", Url: setting.AppSubUrl + "/logout", Icon: "fa fa-fw fa-sign-out", Target: "_self",
			})
		}

		data.NavTree = append(data.NavTree, profileNode)
	}

	if setting.AlertingEnabled && (c.OrgRole == m.ROLE_ADMIN || c.OrgRole == m.ROLE_EDITOR) {
		alertChildNavs := []*dtos.NavLink{
			{Text: "알람 규칙", Id: "alert-list", Url: setting.AppSubUrl + "/alerting/list", Icon: "gicon gicon-alert-rules"},
			{Text: "알림 채널", Id: "channels", Url: setting.AppSubUrl + "/alerting/notifications", Icon: "gicon gicon-alert-notification-channel"},
		}

		data.NavTree = append(data.NavTree, &dtos.NavLink{
			Text:     "알람",
			SubTitle: "알람 규칙 및 알림",
			Id:       "alerting",
			Icon:     "gicon gicon-alert",
			Url:      setting.AppSubUrl + "/alerting/list",
			Children: alertChildNavs,
		})
	}

	enabledPlugins, err := plugins.GetEnabledPlugins(c.OrgId)
	if err != nil {
		return nil, err
	}

	for _, plugin := range enabledPlugins.Apps {
		if plugin.Pinned {
			appLink := &dtos.NavLink{
				Text: plugin.Name,
				Id:   "plugin-page-" + plugin.Id,
				Url:  plugin.DefaultNavUrl,
				Img:  plugin.Info.Logos.Small,
			}

			for _, include := range plugin.Includes {
				if !c.HasUserRole(include.Role) {
					continue
				}

				if include.Type == "page" && include.AddToNav {
					link := &dtos.NavLink{
						Url:  setting.AppSubUrl + "/plugins/" + plugin.Id + "/page/" + include.Slug,
						Text: include.Name,
					}
					appLink.Children = append(appLink.Children, link)
				}

				if include.Type == "dashboard" && include.AddToNav {
					link := &dtos.NavLink{
						Url:  setting.AppSubUrl + "/dashboard/db/" + include.Slug,
						Text: include.Name,
					}
					appLink.Children = append(appLink.Children, link)
				}
			}

			if len(appLink.Children) > 0 && c.OrgRole == m.ROLE_ADMIN {
				appLink.Children = append(appLink.Children, &dtos.NavLink{Divider: true})
				appLink.Children = append(appLink.Children, &dtos.NavLink{Text: "Plugin Config", Icon: "gicon gicon-cog", Url: setting.AppSubUrl + "/plugins/" + plugin.Id + "/edit"})
			}

			if len(appLink.Children) > 0 {
				data.NavTree = append(data.NavTree, appLink)
			}
		}
	}

	if c.IsGrafanaAdmin || c.OrgRole == m.ROLE_ADMIN {
		cfgNode := &dtos.NavLink{
			Id:       "cfg",
			Text:     "설정",
			SubTitle: "기관: " + c.OrgName,
			Icon:     "gicon gicon-cog",
			Url:      setting.AppSubUrl + "/datasources",
			Children: []*dtos.NavLink{
				{
					Text:        "데이터소스",
					Icon:        "gicon gicon-datasources",
					Description: "Add and configure data sources",
					Id:          "datasources",
					Url:         setting.AppSubUrl + "/datasources",
				},
				{
					Text:        "사용자",
					Id:          "users",
					Description: "기관 구성원 관리",
					Icon:        "gicon gicon-user",
					Url:         setting.AppSubUrl + "/org/users",
				},
				{
					Text:        "팀",
					Id:          "teams",
					Description: "기관 그룹 관리",
					Icon:        "gicon gicon-team",
					Url:         setting.AppSubUrl + "/org/teams",
				},
				{
					Text:        "플러그인",
					Id:          "plugins",
					Description: "보기 & 플러그인 관리",
					Icon:        "gicon gicon-plugins",
					Url:         setting.AppSubUrl + "/plugins",
				},
				{
					Text:        "설정",
					Id:          "org-settings",
					Description: "기관 설정",
					Icon:        "gicon gicon-preferences",
					Url:         setting.AppSubUrl + "/org",
				},

				{
					Text:        "API 키",
					Id:          "apikeys",
					Description: "API 키 생성 관리",
					Icon:        "gicon gicon-apikeys",
					Url:         setting.AppSubUrl + "/org/apikeys",
				},
			},
		}

		if c.OrgRole != m.ROLE_ADMIN {
			cfgNode = &dtos.NavLink{
				Id:       "cfg",
				Text:     "설정",
				SubTitle: "기관: " + c.OrgName,
				Icon:     "gicon gicon-cog",
				Url:      setting.AppSubUrl + "/admin/users",
				Children: make([]*dtos.NavLink, 0),
			}
		}

		if c.OrgRole == m.ROLE_ADMIN && c.IsGrafanaAdmin {
			cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
				Divider: true, HideFromTabs: true, Id: "admin-divider", Text: "Text",
			})
		}

		if c.IsGrafanaAdmin {
			cfgNode.Children = append(cfgNode.Children, &dtos.NavLink{
				Text:         "서버 관리",
				HideFromTabs: true,
				SubTitle:     "모든 기관과 사용자 관리",
				Id:           "admin",
				Icon:         "gicon gicon-shield",
				Url:          setting.AppSubUrl + "/admin/users",
				Children: []*dtos.NavLink{
					{Text: "사용자", Id: "global-users", Url: setting.AppSubUrl + "/admin/users", Icon: "gicon gicon-user"},
					{Text: "기관", Id: "global-orgs", Url: setting.AppSubUrl + "/admin/orgs", Icon: "gicon gicon-org"},
					{Text: "설정", Id: "server-settings", Url: setting.AppSubUrl + "/admin/settings", Icon: "gicon gicon-preferences"},
					{Text: "통계", Id: "server-stats", Url: setting.AppSubUrl + "/admin/stats", Icon: "fa fa-fw fa-bar-chart"},
					{Text: "스타일가이드", Id: "styleguide", Url: setting.AppSubUrl + "/styleguide", Icon: "fa fa-fw fa-eyedropper"},
				},
			})
		}

		data.NavTree = append(data.NavTree, cfgNode)
	}

	data.NavTree = append(data.NavTree, &dtos.NavLink{
		Text:         "도움",
		SubTitle:     fmt.Sprintf(`%s v%s (%s)`, setting.ApplicationName, setting.BuildVersion, setting.BuildCommit),
		Id:           "help",
		Url:          "#",
		Icon:         "gicon gicon-question",
		HideFromMenu: true,
		Children: []*dtos.NavLink{
			{Text: "단축키", Url: "/shortcuts", Icon: "fa fa-fw fa-keyboard-o", Target: "_self"},
		},
	})

	return &data, nil
}

func Index(c *m.ReqContext) {
	data, err := setIndexViewData(c)
	if err != nil {
		c.Handle(500, "Failed to get settings", err)
		return
	}
	c.HTML(200, "index", data)
}

func NotFoundHandler(c *m.ReqContext) {
	if c.IsApiRequest() {
		c.JsonApiErr(404, "Not found", nil)
		return
	}

	data, err := setIndexViewData(c)
	if err != nil {
		c.Handle(500, "Failed to get settings", err)
		return
	}

	c.HTML(404, "index", data)
}
