import StyleSheet from 'react-style';
import React from 'react';
import moment from 'moment';

import {User} from '../models/Models';
import CardFormMixin from './CardFormMixin';
import AdminAuthorizeMixin from './AdminAuthorizeMixin.js'

let EditSecretCard = React.createClass({
  mixins: [CardFormMixin],

  render() {
    const app = this.props.app;
    const proc = this.props.proc;
    const path = this.props.path;
    let content = this.props.content;
    let rows = 5 + this.getRows(content);
    return (
      <div className="mdl-card mdl-shadow--2dp" styles={[this.styles.card, this.props.style]}>
        <div className="mdl-card__title">
          <h2 className="mdl-card__title-text">
            { `编辑 ${app}, ${proc} 的 ${path} 的内容`}
          </h2>
        </div>
        <form style={this.styles.form}>
          <div className="mdl-textfield mdl-js-textfield mdl-textfield--floating-label" style={{ width: '100%' }}>
            <textarea className="mdl-textfield__input" type="text" id="text" ref="content" rows={rows}>
              {content}
            </textarea>
          </div>
        </form>
        <div className="mdl-card__actions" style={{ textAlign: 'right'  }}>
          { this.renderActions("备份并更新", this.onBackupAndUpdate) }
          { this.renderActions("更新", this.onUpdate)}
        </div>
      </div>
    );
  },

  getRows(content){
    let num = 0;
    for(var i=0; i < content.length-1; i++){
      if (content.charAt(i) =='\n'){
        num = num + 1;
      }
    }
    return num;
  },

  onBackupAndUpdate() {
    const app = this.props.app;
    const proc = this.props.proc;
    const path = this.props.path;
    const oldContent = this.props.content;
    let node = this.refs['content'].getDOMNode();
    let oldFormData = {};
    oldFormData['appname']=app;
    oldFormData['procname']=proc;
    oldFormData['fpath']=path+'-'+Math.round(moment.now()/1000);
    oldFormData['content'] = oldContent;
    let newFormData = {};
    newFormData['appname']=app;
    newFormData['procname']=proc;
    newFormData['fpath']=path;
    newFormData['content']=node.value;
    this.setState({ inRequest: true  });
    User.putSecret(oldFormData, this.onRequestCallback);
    User.putSecret(newFormData, this.onRequestCallback);
  },

  onUpdate(){
    const app = this.props.app;
    const proc = this.props.proc;
    const path = this.props.path;
    let node = this.refs['content'].getDOMNode();
    let newFormData = {};
    newFormData['appname']=app;
    newFormData['procname']=proc;
    newFormData['fpath']=path;
    newFormData['content']=node.value;
    this.setState({ inRequest: true  });
    User.putSecret(newFormData, this.onRequestCallback);
  },

});

export default EditSecretCard;
