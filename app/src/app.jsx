import * as React from 'react'
import {useRoutes} from 'hookrouter'

import Home from './home'
import Signature from './signature'
import Template from './template'
import NotFound from './not-found'
import Layout from './layout'

import 'normalize.css'

const routes = {
  '/': () => <Home />,
  '/signature/:username/:id/:mac': ({ username id, mac }) => <Signature id={id} username={username} mac={mac} />,
  '/template/:id': ({ id }) => <Template id={id} />,
}

const App = () => {
  const routeResult = useRoutes(routes)

  return (
    <Layout>{routeResult || <NotFound />}</Layout>
  )
}

export default App
