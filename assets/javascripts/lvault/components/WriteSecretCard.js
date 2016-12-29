import StyleSheet from 'react-style';
import React from 'react';

import {User} from '../models/Models';
import CardFormMixin from './CardFormMixin';
import AdminAuthorizeMixin from './AdminAuthorizeMixin.js'

let WriteSecretCard = React.createClass({
  mixins: [CardFormMixin],

  getInitialState() {
    return {
      formValids: {
        'appname': true,
        'procname': true,
		'fpath': true,
	  },
    };
  },

  render() {
    const {reqResult} = this.state;
    return (
      <div className="mdl-card mdl-shadow--2dp" styles={[this.styles.card, this.props.style]}>
        <div className="mdl-card__title">
			<h2 className="mdl-card__title-text">写入秘密文件</h2>
        </div>
        { this.renderResult() }
        { 
          reqResult.fin && reqResult.ok ? null :
            this.renderForm(null, [
              this.renderInput("appname", "APP name*(字母、数字、减号和'.')", { type: "text", pattern: "[\-a-zA-Z0-9.]*" }),
              this.renderInput("procname", "proc 全名*(eg appname.web.procname)", { type: 'text' }),
			  this.renderInput("fpath", "绝对路径*(eg /lain/app/hello.dat)", { type: 'text' }),
              this.renderTextArea("content", "当前 secret file 的内容", { type: 'text' }), 
            ])
        }
        { this.renderAction("提交", this.onWriteSecret) }
      </div>
    );
  },

  onWriteSecret() {
	const fields = ['appname', 'procname', 'content', 'fpath'];
    const rFields = ['appname', 'procname', 'fpath'];
    const {isValid, formData} = this.validateForm(fields, rFields);
    if (isValid) {
      this.setState({ inRequest: true });
      User.putSecret(formData, this.onRequestCallback);
    }
  },

});

export default WriteSecretCard;
