"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var core_1 = require('@angular/core');
var platform_browser_1 = require('@angular/platform-browser');
// import { ContentComponent }  from './content.component';
// import { NavComponent } from './nav.component';
// import { RepoListComponent }  from './repoList.component';
// import { RepositoriesComponent }  from './repositories.component';
// import { RepoCreateComponent }  from './repoCreate.component';
// import { RepoDetailComponent } from './repoDetail.component';
// import { OrgListComponent }      from './orgList.component';
// import { OrgCreateComponent }      from './orgCreate.component';
// import { OrgEditComponent }      from './orgEdit.component';
// import { UserSettingComponent }      from './userSetting.component';
// import { PromptComponent }  from './prompt.component';
// import { PromptModule }  from './prompt.module';
// import { OrgService }      from './org.service';
// import { RepoService }      from './repo.service';
var forms_1 = require('@angular/forms');
var http_1 = require('@angular/http');
// import { routing } from './content.routing';
var IndexModule = (function () {
    function IndexModule() {
    }
    IndexModule = __decorate([
        core_1.NgModule({
            imports: [
                platform_browser_1.BrowserModule,
                forms_1.FormsModule,
                http_1.HttpModule,
            ],
            declarations: [],
            providers: [],
        }), 
        __metadata('design:paramtypes', [])
    ], IndexModule);
    return IndexModule;
}());
exports.IndexModule = IndexModule;
