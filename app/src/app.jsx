import * as React from 'react'
import {useRoutes} from 'hookrouter'
import { pdfjs } from 'react-pdf'

import Home from './home'
import Signature from './signature'
import Document from './document'
import NotFound from './not-found'
import Layout from './layout'

import 'normalize.css'

pdfjs.GlobalWorkerOptions.workerSrc = `//cdnjs.cloudflare.com/ajax/libs/pdf.js/${pdfjs.version}/pdf.worker.js`

const routes = {
  '/': () => <Home />,
  '/signature/:username/:id/:hash': ({ username, id, hash }) => <Signature id={id} username={username} hash={hash} />,
  '/document/:username/:id/:hash': ({ username, id, hash }) => <Document id={id} username={username} hash={hash} />,
}

const App = () => {
  const routeResult = useRoutes(routes)

  return (
    <Layout>{routeResult || <NotFound />}</Layout>
  )
}

export default App
