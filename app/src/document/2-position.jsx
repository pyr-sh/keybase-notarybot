import * as React from 'react'
import clsx from 'clsx'
import useWindowSize from '../utils/use-window-size'
import useBodyID from '../utils/use-body-id'
import { useDrag, useWheel } from 'react-use-gesture'
import { Document, Page } from 'react-pdf'
// import { Document, Page } from 'react-pdf/dist/esm/entry.webpack'

const Position = ({ document }) => {
  useBodyID('document-position')

  const size = useWindowSize()

  const [page, setPage] = React.useState(1)
  const [zoom, setZoom] = React.useState(1)
  const [posX, setPosX] = React.useState(0)
  const [posY, setPosY] = React.useState(0)

  const [numPages, setNumPages] = React.useState(0)
  const handleLoadSuccess = React.useCallback(({ numPages }) => {
    setNumPages(numPages)
  }, [setNumPages])

  const wheelHandler = React.useCallback(({delta: [, y]}) => {
    const step = (y / 53) * 0.1
    setZoom(zoom - step < 0.25 ? 0.25 : zoom - step)
  }, [zoom, setZoom])
  const wheelBind = useWheel(wheelHandler, {domTarget: window})

  const dragHandler = React.useCallback(({delta: [x, y]}) => {
    setPosX(posX + x)
    setPosY(posY + y)
  }, [posX, setPosX, posY, setPosY])
  const dragBind = useDrag(dragHandler, {domTarget: window})

  const [signatories, setSignatories] = React.useState([
    {name: 'Person #1'},
    {name: 'Person #2'},
  ])

  return (
    <div className="document-position" {...dragBind()}>
      <div className="document-document" {...wheelBind()} style={{
        left: posX,
        top: posY,
      }}>
        <Document file={document} onLoadSuccess={handleLoadSuccess}>
          {
            numPages > 0 &&
              <Page
                renderAnnotationLayer={false}
                renderTextLayer={false}
                pageNumber={page}
                height={size.height * zoom}
              />
          }
        </Document>
      </div>

      <div className="document-zoom">{Math.floor(zoom * 100)}%</div>
      <div className="document-page">{page} / {numPages}</div>
      <div className="document-palette">
        <h3>Signatories</h3>
        {signatories.map((person, i) => (
          <div key={i} className={clsx('document-signatory', {'document-signatory-placed': person.x !== undefined})}>
            {person.name}
          </div>
        ))}
      </div>
    </div>
  )
}

export default Position
