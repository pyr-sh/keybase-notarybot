import * as React from 'react'

import '@openfonts/open-sans_all'
import './style.css'

type Props = {
  children: React.ReactNode
}

const Layout = (props: Props) => {
  return (
    <div id="root">
      <div id="layout-logo">
        <img alt="logo" src="https://s3.amazonaws.com/keybase_processed_uploads/ec9e9600d9dc4b940dab0e1e1fdcb705_360_360.jpg" />
        <span>notarybot</span>
      </div>

      {props.children}
    </div>
  )
}

export default Layout
