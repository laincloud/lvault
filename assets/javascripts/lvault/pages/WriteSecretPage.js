import StyleSheet from 'react-style';
import React from 'react';

import WriteSecretCard from '../components/WriteSecretCard';

let WriteSecretPage = React.createClass({
  render() {
    return (
      <div className="mdl-grid">
        <div className="mdl-cell mdl-cell--6-col mdl-cell--8-col-tablet mdl-cell--4-col-phone">
			<WriteSecretCard />
        </div>
      </div>
    );
  }
});

export default WriteSecretPage;
