import StyleSheet from 'react-style';
import React from 'react';
import {History} from 'react-router';

import {Admin} from '../models/Models';
import {User} from '../models/Models';

let ListSecretsCard = React.createClass({
  mixins: [History],

  getInitialState(){
    return {
      secrets: this.getSecrets(),
    };
  },

  componentDidMount() {
    this.reload(); 
  },

  render() {
    let secrets = JSON.parse(this.state.secrets);
    secrets = this.sortAndGroup(secrets);
    return (
      <div className="mdl-cell mdl-cell--12-col mdl-cell--12-col-tablet mdl-cell--4-col-phone">
        {
          secrets.map((secret,index) => {
            return (
              this.renderProc(secret)
            );
          })
        }
     </div>
    );
  },

  renderProc(secret){
    return (
       <div className="mdl-card mdl-shadow--2dp" styles={[this.styles.card, this.props.style]}>
         
         <div className="mdl-card__title">
          <h2 className="mdl-card__title-text"> 
            {`Proc - ${secret.procname}: secret files 列表`}
          </h2>
        </div>

        <table className="mdl-data-table mdl-js-data-table" style={this.styles.table}>
          <thead>
            <tr>
			  <th className="mdl-data-table__cell--non-numeric">路径</th>
              <th className="mdl-data-table__cell--non-numeric">内容</th>
              <th className="mdl-data-table__cell--non-numeric"></th>
              <th className="mdl-data-table__cell--non-numeric"></th>
            </tr>
          </thead>
          <tbody>
            {
              secret.secrets.map((s, index) => {
                return (
                  <tr key={`s-${index}`}>
                     <td className="mdl-data-table__cell--non-numeric">{s.path}</td>
                     <td className="mdl-data-table__cell--non-numeric" style={this.styles.breaklineTd}>{s.content}</td>
                     <td className="mdl-data-table__cell--non-numeric">
                        <a href="javascript:;" key={`edit-${s.path}`} 
                          style={{ marginRight: 8 }}
                          onClick={(evt) => this.goEditSecretFile(this.props.app,secret.procname, s.path, s.content)}>编辑</a>
                     </td>
                     <td className="mdl-data-table__cell--non-numeric">
                        <a href="javascript:;" key={`edit-${s.path}`} 
                          style={{ marginRight: 8 }}
                          onClick={(evt) => this.goDeleteSecretFile(this.props.app,secret.procname, s.path)}>删除</a>
                     </td>
 
                  </tr>
                  );
              })
            }
          </tbody>
        </table>

      </div>
    );
  },

  getSecrets(){
    let secrets = this.props.secrets;
    if(secrets == undefined){
      secrets = "";
    }
    return secrets;
  },

  goDeleteSecretFile(app,proc,path){
    let yes = confirm(`确定要删除文件 - ${path} 吗？`);
    if (yes){
      let formData = {};
      formData['appname']=app;
      formData['procname']=proc;
      formData['fpath']=path
      User.deleteSecret(formData, this.reload);
    }
  },

  reload(){
    let formData={};
    formData['appname']=this.props.app;
    formData['procname']=this.props.proc;
    User.getSecrets(formData, this.updateState);
  },

  updateState(ok, status){
    if(!ok){
      console.log(status);
    }else{
      let secrets = status;
      this.setState({ secrets });
    }
  },

  goEditSecretFile(app,proc,path,content){
    this.history.pushState({
      app,
      proc,
      path,
      content,
    }, `/v2/spa/secret-edit/${app}`)
  },

  getProcName(path){
    let from = path.indexOf("/");
    let next = path.indexOf("/", from+1);
    return path.substr(from + 1, next-from-1);
  },

  getFilePath(path){
    let from = path.indexOf("/");
    let next = path.indexOf("/", from+1);
    let fileBegin = path.indexOf("/", next+1);
    return path.substr(next, path.length - next);
  },

  sortAndGroup(secrets){
    var procMap = new Map();
    for(var i=0; i < secrets.length; i++){
      let procname = this.getProcName(secrets[i]['path']);
      procMap.set(procname,'')
    }
    var procs = []
    for(var key of procMap.keys()){
      procs.push(key)
    }
    procs.sort();
    var ret = [];
    for(var j=0; j< procs.length; j++){
      ret.push({})
    }
    for(var j=0; j< procs.length; j++){
      ret[j]['procname']=procs[j]
    }
    for(var i=0; i< secrets.length; i++){
      let procname = this.getProcName(secrets[i]['path']);
      for(var j=0; j< procs.length; j++){
        if(procs[j] == procname){
          let filePath = this.getFilePath(secrets[i]['path']);
          if(ret[j]['secrets']==undefined){
            ret[j]['secrets']=[];
          }
          ret[j]['secrets'].push({'path':filePath, 'content':secrets[i]['content']})
        }
      }
    }
    return ret
  },

  styles: StyleSheet.create({
    card: {
      width: '100%',
      marginBottom: 16,
      minHeight: 50,
    },
    table: {
      width: '100%',
      borderLeft: 'none',
      borderRight: 'none',
    },
    breaklineTd: {
      whiteSpace: 'pre-line',
      wordWrap: 'break-word',
      wordBreak: 'break-word',
    },
  }),
});

export default ListSecretsCard;
