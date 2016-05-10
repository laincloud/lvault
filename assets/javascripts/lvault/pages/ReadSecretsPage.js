import StyleSheet from 'react-style';
import React from 'react';

import ReadSecretsCard from '../components/ReadSecretsCard';

let ReadSecretsPage = React.createClass({
  render() {
    return (
      <div className="mdl-grid">
        <div className="mdl-cell mdl-cell--6-col mdl-cell--8-col-tablet mdl-cell--4-col-phone">
			<ReadSecretsCard />
        </div>
      </div>
    );
  }
});

export default ReadSecretsPage;
