import StyleSheet from 'react-style';
import React from 'react';

import ListSecretsCard from '../components/ListSecretsCard';

let ListSecretsPage = React.createClass({

  getInitialState(){
    return{
      params:{
        'app':this.props.params.app,
        'proc':this.props.params.proc,
        'secrets': this.props.location.state.status,
      },
    };
  },

  render() {
    const {app,proc,secrets}=this.state.params;
    return (
      <div className="mdl-grid">
		<ListSecretsCard app={app} proc={proc} secrets={secrets}/>
      </div>
    );
  }
});

export default ListSecretsPage;
