/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import { Component, ElementRef, ViewChild, OnInit } from '@angular/core';
import { Http, Response } from '@angular/http';
import { RepoService } from './repo.service';
import { FormBuilder, FormGroup, Validators, FormControl } from '@angular/forms';
import { CustomValidators } from 'ng2-validation';

const password = new FormControl('', Validators.required);
const confirmPassword = new FormControl('', CustomValidators.equalTo(password));

//import { Repo } from './repo';
// import * as D3 from 'd3';
// import * as yaml from 'js-yaml';

export class Repo {
  name: string;
  repo_type: number;
  islimit:boolean;
  description:string
  ishavekey:boolean;
}

export class DockerIMG {
  name: string;
  description:string
  cmd:string
  path: string;
  tag:string;
  version:string
  isCer:boolean;
}
@Component({
  selector: 'app-repo',
  templateUrl: './repo.component.html',
  styleUrls: ['./repo.component.scss']
})


export class RepoComponent implements OnInit {
//  @ViewChild('design') element: ElementRef;
 // constructor(private http: Http) {};
  
  public form: FormGroup;
  constructor(private fb: FormBuilder) {}

  // control: FormControl = new FormControl('value', Validators.minLength(2));
  // setValue() { this.control.setValue('new value'); }

    username ="DeanLee"
    webside_url ="hubops-docker-test_private.bintray.io"
    isCreate = false
    repo: Repo  = {
      name: '',
      islimit:true,
      repo_type:-1,
      ishavekey:false,
      description:""
    };

    // dockerimg: DockerIMG 

    selecteddocker:DockerIMG;
    
    onSelect(dockerimg: DockerIMG): void {
        dockerimg =  {
        name: this.repo.name+"_dockerimg",
        description:"",
        cmd:"string",
        path: "string",
        tag:"string",
        version:"string",
        isCer:false
      }
      this.selecteddocker = dockerimg;
    }

//    constructor(private repoService: RepoService) { }

//    getrepo(): void {
//     this.repoService.addRepo().then(repo => this.repo = repo);
//   }

    onCreateButtonClick(repo: Repo): void {
     this.isCreate =true;
      alert(JSON.stringify( repo))
    }
 
        
    ngOnInit() {
    this.form = this.fb.group({
          fname: [null, Validators.compose([Validators.required, Validators.minLength(5), Validators.maxLength(10)])],
          desc: [null, Validators.compose([Validators.minLength(5), Validators.maxLength(100)])]
        });
        
      }
}

