import { Injectable } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { Http, Headers, Response, RequestOptions } from '@angular/http';
import 'rxjs/add/operator/toPromise';
import 'rxjs/add/operator/catch';


@Injectable()
export class UserService {

  private headers = new Headers({'Content-Type': 'application/json'});
  constructor(
    private http: Http,
    private title: Title
  ){}
  
  getBrowseList(): Promise<any> {
    return this.http.get('json/browseList.json')
               .toPromise()
               .then(this.dealData,this.dealError)
               .catch(this.handleError)
  }

  doLogin(info): Promise<any> {
    let params=JSON.stringify(info)
    return this.http.post('/web/v1/user/signin', params, {headers: this.headers})
               .toPromise()
               .then(this.dealData,this.dealError)
               .catch(this.handleError)
  }

  signUp(info): Promise<any> {
    let params=JSON.stringify(info)
    return this.http.post('/web/v1/user', params, {headers: this.headers})
               .toPromise()
               .then(this.dealData,this.dealError)
               .catch(this.handleError)
  }

  sendEmail(info): Promise<any> {
    let params=JSON.stringify(info)
    return this.http.post('/web/v1/user/forget', params, {headers: this.headers})
               .toPromise()
               .then(this.dealData,this.dealError)
               .catch(this.handleError)
  }

  resetPwd(info): Promise<any> {
    let params=JSON.stringify(info)
    return this.http.post('/web/v1/user/forget/reset', params, {headers: this.headers})
               .toPromise()
               .then(this.dealData,this.dealError)
               .catch(this.handleError)
  }

  getEmailList(user): Promise<any> {
    return this.http.get('/web/v1/user/'+user.username+'/emails')
               .toPromise()
               .then(this.dealData,this.dealError)
               .catch(this.handleError)
  }

  addEmail(info): Promise<any> {
    let params=JSON.stringify(info)
    return this.http.put('/web/v1/user/'+info.username+'/email', params, {headers: this.headers})
               .toPromise()
               .then(this.dealData,this.dealError)
               .catch(this.handleError)
  }

  verifyEmail(info,user): Promise<any> {
    info.username = user.username;
    let params=JSON.stringify(info)
    return this.http.put('/web/v1/user/'+user.username+'/email/'+info.id+'/send', params, {headers: this.headers})
               .toPromise()
               .then(this.dealData)
               .catch(this.handleError)
  }
  
  delEmail(info,user): Promise<any> {
    let params=JSON.stringify(info)
    return this.http.delete('/web/v1/user/'+user.username+'/email/'+info.id, {headers: this.headers})
               .toPromise()
               .then(this.dealData,this.dealError)
               .catch(this.handleError)
  }

  loginOut(user): Promise<any> {
    let params=JSON.stringify(user)
    return this.http.put('/web/v1/user/'+user.username+'/signout', params, {headers: this.headers})
               .toPromise()
               .then(this.dealData,this.dealError)
               .catch(this.handleError)
  }

  private dealData (res: Response) {
    var object = {
      code: res.status,
      data: res.json()
    }
    console.log(res)
    return object || {}
  }

  private dealError (err: Response) {
    var object = {
      code: err.status,
      data: err.json()
    }
    console.log(err)
    return object || {}
  }

  private handleError (error: any) {
    console.log(error)
    // let errMsg = (error.message) ? error.message :
    //   error.status ? `${error.status} - ${error.statusText}` : 'Server error';
    // console.log(errMsg);
    var object = {
      code: error.status,
      data: error.json()
    }; 
    return Promise.reject(object);
  }

  changeTitle(val){
    this.title.setTitle(val)
  }
}