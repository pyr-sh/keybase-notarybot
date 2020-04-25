import * as React from 'react'
import { Document, Page } from 'react-pdf'
// import { Document, Page } from 'react-pdf/dist/esm/entry.webpack'

const Position = ({ document }) => {
  const [numPages, setNumPages] = React.useState(0)
  const handleLoadSuccess = React.useCallback(({ numPages }) => {
    setNumPages(numPages)
  }, [setNumPages])

  return (
    <div className="document-position">
      <Document file={document} onLoadSuccess={handleLoadSuccess}>
        {
          Array.from(
            new Array(numPages),
            (el, index) => (
              <Page
                key={`page_${index + 1}`}
                pageNumber={index + 1}
              />
            ),
          )
        }
      </Document>
    </div>
  )
}

export default Position
