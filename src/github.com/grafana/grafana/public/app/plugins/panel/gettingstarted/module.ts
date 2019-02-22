import { PanelCtrl } from 'app/plugins/sdk';

import { contextSrv } from 'app/core/core';

class GettingStartedPanelCtrl extends PanelCtrl {
  static templateUrl = 'public/app/plugins/panel/gettingstarted/module.html';
  checksDone: boolean;
  stepIndex: number;
  steps: any;

  /** @ngInject **/
  constructor($scope, $injector, private backendSrv, datasourceSrv, private $q) {
    super($scope, $injector);

    this.stepIndex = 0;
    this.steps = [];

    this.steps.push({
      title: '아르고스 데이터 시각화 엔진 설치',
      icon: 'icon-gf icon-gf-check',
      href: '#',
      target: '_blank',
      note: '설치 문서 확인',
      check: () => $q.when(true),
    });

    this.steps.push({
      title: '첫 데이터 소스를 생성하세요',
      cta: '데이터 소스 추가',
      icon: 'icon-gf icon-gf-datasources',
      href: 'datasources/new?gettingstarted',
      check: () => {
        return $q.when(
          datasourceSrv.getMetricSources().filter(item => {
            return item.meta.builtIn !== true;
          }).length > 0
        );
      },
    });

    this.steps.push({
      title: '첫 대쉬보드를 생성하세요',
      cta: '새 대쉬보드',
      icon: 'icon-gf icon-gf-dashboard',
      href: 'dashboard/new?gettingstarted',
      check: () => {
        return this.backendSrv.search({ limit: 1 }).then(result => {
          return result.length > 0;
        });
      },
    });

    this.steps.push({
      title: '팀 초대',
      cta: '사용자 추가',
      icon: 'icon-gf icon-gf-users',
      href: 'org/users?gettingstarted',
      check: () => {
        return this.backendSrv.get('/api/org/users').then(res => {
          return res.length > 1;
        });
      },
    });

    this.steps.push({
      title: '플러그인 및 어플리케이션 추가',
      cta: '플러그인 리포지토리 검색',
      icon: 'icon-gf icon-gf-apps',
      href: 'https://argos.goback.world/plugins?utm_source=grafana_getting_started',
      check: () => {
        return this.backendSrv.get('/api/plugins', { embedded: 0, core: 0 }).then(plugins => {
          return plugins.length > 0;
        });
      },
    });
  }

  $onInit() {
    this.stepIndex = -1;
    return this.nextStep().then(res => {
      this.checksDone = true;
    });
  }

  nextStep() {
    if (this.stepIndex === this.steps.length - 1) {
      return this.$q.when();
    }

    this.stepIndex += 1;
    var currentStep = this.steps[this.stepIndex];
    return currentStep.check().then(passed => {
      if (passed) {
        currentStep.cssClass = 'completed';
        return this.nextStep();
      }

      currentStep.cssClass = 'active';
      return this.$q.when();
    });
  }

  dismiss() {
    this.dashboard.removePanel(this.panel, false);

    this.backendSrv
      .request({
        method: 'PUT',
        url: '/api/user/helpflags/1',
        showSuccessAlert: false,
      })
      .then(res => {
        contextSrv.user.helpFlags1 = res.helpFlags1;
      });
  }
}

export { GettingStartedPanelCtrl, GettingStartedPanelCtrl as PanelCtrl };