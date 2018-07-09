import coreModule from 'app/core/core_module';

const hitTypes = {
  FOLDER: 'dash-folder',
  DASHBOARD: 'dash-db',
};

export class ValidationSrv {
  rootName = 'general';

  /** @ngInject */
  constructor(private $q, private backendSrv) {}

  validateNewDashboardName(folderId, name) {
    return this.validate(folderId, name, '이 폴더안에 동일한 이름의 대쉬보드가 있습니다');
  }

  validateNewFolderName(name) {
    return this.validate(0, name, '동일한 이름의 대쉬보드가 이미 있습니다.');
  }

  private validate(folderId, name, existingErrorMessage) {
    name = (name || '').trim();
    const nameLowerCased = name.toLowerCase();

    if (name.length === 0) {
      return this.$q.reject({
        type: 'REQUIRED',
        message: '이름 필수',
      });
    }

    if (folderId === 0 && nameLowerCased === this.rootName) {
      return this.$q.reject({
        type: 'EXISTING',
        message: '예약어이거나 폴더 이름으로 사용할 수 없습니다.',
      });
    }

    let deferred = this.$q.defer();

    const promises = [];
    promises.push(this.backendSrv.search({ type: hitTypes.FOLDER, folderIds: [folderId], query: name }));
    promises.push(this.backendSrv.search({ type: hitTypes.DASHBOARD, folderIds: [folderId], query: name }));

    this.$q.all(promises).then(res => {
      let hits = [];

      if (res.length > 0 && res[0].length > 0) {
        hits = res[0];
      }

      if (res.length > 1 && res[1].length > 0) {
        hits = hits.concat(res[1]);
      }

      for (let hit of hits) {
        if (nameLowerCased === hit.title.toLowerCase()) {
          deferred.reject({
            type: 'EXISTING',
            message: existingErrorMessage,
          });
          break;
        }
      }

      deferred.resolve();
    });

    return deferred.promise;
  }
}

coreModule.service('validationSrv', ValidationSrv);
