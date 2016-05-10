import StyleSheet from 'react-style';
import React from 'react';
import {History} from 'react-router';

import AdminDeleteSecretCard from '../components/AdminDeleteSecretCard';
import AdminAuthorizeMixin from '../components/AdminAuthorizeMixin';
import {Admin} from '../models/Models';

let AdminUsersPage = React.createClass({
  mixins: [History, AdminAuthorizeMixin],

  componentWillMount() {
	   this.authorize('delete');
  },

  render() {
	   const isValid = this.isSessionValid();
    return (
      <div className="mdl-grid">
        <div className="mdl-cell mdl-cell--6-col mdl-cell--8-col-tablet mdl-cell--4-col-phone">
          { 
			  !isValid ? <p>等待认证中....</p> :
			  <AdminDeleteSecretCard token={this.state.token} tokenType={this.state.tokenType} />
          }
        </div>
      </div>
    );
  },

  styles: StyleSheet.create({
  }),

});

export default AdminUsersPage;
