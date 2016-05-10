import StyleSheet from 'react-style';
import React from 'react';
import {History} from 'react-router';

import AdminListSynctimeCard from '../components/AdminListSynctimeCard';
import AdminAuthorizeMixin from '../components/AdminAuthorizeMixin';
import {Admin} from '../models/Models';

let AdminSynctimePage = React.createClass({
  mixins: [History, AdminAuthorizeMixin],

  componentWillMount() {
    this.authorize('syncstatus');
  },

  render() {
    const isValid = this.isSessionValid();
    return (
      <div className="mdl-grid">
        <div className="mdl-cell mdl-cell--12-col mdl-cell--8-col-tablet mdl-cell--4-col-phone">
          {
            !isValid ? <p>等待认证中……</p>
              : <AdminListSynctimeCard ref="groupList" token={this.state.token} tokenType={this.state.tokenType} />
          }
        </div>
      </div>
    );
  },

  refreshGroups() {
    if (this.isSessionValid()) {
      this.refs.groupList.reload();
    }
  },

  styles: StyleSheet.create({
  }),

});

export default AdminSynctimePage;
