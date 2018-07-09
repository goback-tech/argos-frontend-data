import _ from 'lodash';
import { types, getEnv } from 'mobx-state-tree';
import { NavItem } from './NavItem';

export const NavStore = types
  .model('NavStore', {
    main: types.maybe(NavItem),
    node: types.maybe(NavItem),
  })
  .actions(self => ({
    load(...args) {
      let children = getEnv(self).navTree;
      let main, node;
      let parents = [];

      for (let id of args) {
        node = children.find(el => el.id === id);

        if (!node) {
          throw new Error(`NavItem with id ${id} not found`);
        }

        children = node.children;
        parents.push(node);
      }

      main = parents[parents.length - 2];

      if (main.children) {
        for (let item of main.children) {
          item.active = false;

          if (item.url === node.url) {
            item.active = true;
          }
        }
      }

      self.main = NavItem.create(main);
      self.node = NavItem.create(node);
    },

    initFolderNav(folder: any, activeChildId: string) {
      let main = {
        icon: 'fa fa-folder-open',
        id: 'manage-folder',
        subTitle: '폴더 대쉬보드와 권한 관리',
        url: '',
        text: folder.title,
        breadcrumbs: [{ title: '대쉬보드', url: 'dashboards' }],
        children: [
          {
            active: activeChildId === 'manage-folder-dashboards',
            icon: 'fa fa-fw fa-th-large',
            id: 'manage-folder-dashboards',
            text: '대쉬보드',
            url: folder.url,
          },
          {
            active: activeChildId === 'manage-folder-permissions',
            icon: 'fa fa-fw fa-lock',
            id: 'manage-folder-permissions',
            text: '권한',
            url: `${folder.url}/permissions`,
          },
          {
            active: activeChildId === 'manage-folder-settings',
            icon: 'fa fa-fw fa-cog',
            id: 'manage-folder-settings',
            text: '설정',
            url: `${folder.url}/settings`,
          },
        ],
      };

      self.main = NavItem.create(main);
    },

    initDatasourceEditNav(ds: any, plugin: any, currentPage: string) {
      let title = '신규';
      let subTitle = `Type: ${plugin.name}`;

      if (ds.id) {
        title = ds.name;
      }

      let main = {
        img: plugin.info.logos.large,
        id: 'ds-edit-' + plugin.id,
        subTitle: subTitle,
        url: '',
        text: title,
        breadcrumbs: [{ title: '데이터소스', url: 'datasources' }],
        children: [
          {
            active: currentPage === 'datasource-settings',
            icon: 'fa fa-fw fa-sliders',
            id: 'datasource-settings',
            text: '설정',
            url: `datasources/edit/${ds.id}`,
          },
        ],
      };

      const hasDashboards = _.find(plugin.includes, { type: 'dashboard' }) !== undefined;
      if (hasDashboards && ds.id) {
        main.children.push({
          active: currentPage === 'datasource-dashboards',
          icon: 'fa fa-fw fa-th-large',
          id: 'datasource-dashboards',
          text: '대쉬보드',
          url: `datasources/edit/${ds.id}/dashboards`,
        });
      }

      self.main = NavItem.create(main);
    },
  }));
