import StyleSheet from 'react-style';
import React from 'react';

import EditSecretCard from '../components/EditSecretCard';

let EditSecretPage = React.createClass({

  render() {
    console.log(this.props);
    const app = this.props.location.state.app;
    const proc = this.props.location.state.proc;
    const path = this.props.location.state.path;
    let content = this.props.location.state.content;
    return (
      <div className="mdl-grid">
        <div className="mdl-cell mdl-cell--12-col mdl-cell--12-col-tablet mdl-cell--4-col-phone">
			<EditSecretCard app={app} proc={proc} path={path} content={content}/>
        </div>
      </div>
    );
  }
});

export default EditSecretPage;
