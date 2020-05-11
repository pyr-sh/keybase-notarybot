import * as React from "react";
import clsx from "clsx";
import useWindowSize from "../utils/use-window-size";
import useBodyID from "../utils/use-body-id";
import { useDrag, useWheel } from "react-use-gesture";
import { Document, Page } from "react-pdf";
import ReactCrop from "react-image-crop";
// import { Document, Page } from 'react-pdf/dist/esm/entry.webpack'

const Position = ({ document, onBack, onSave }) => {
  useBodyID("document-position");

  const size = useWindowSize();

  const [page, setPage] = React.useState(1);
  const [zoom, setZoom] = React.useState(1);
  const [posX, setPosX] = React.useState(0);
  const [posY, setPosY] = React.useState(0);

  const [numPages, setNumPages] = React.useState(0);
  const handleLoadSuccess = React.useCallback(
    ({ numPages }) => {
      setNumPages(numPages);
    },
    [setNumPages]
  );
  const handlePreviousPage = React.useCallback(() => {
    if (numPages === 0) {
      return;
    }
    if (page === 1) {
      return;
    }
    setPage(page - 1);
  }, [page, setPage, numPages]);
  const handleNextPage = React.useCallback(() => {
    if (numPages === 0) {
      return;
    }
    if (page === numPages) {
      return;
    }
    setPage(page + 1);
  }, [page, setPage, numPages]);
  const handleBack = React.useCallback(() => onBack(), [onBack]);

  const wheelHandler = React.useCallback(
    ({ delta: [, y] }) => {
      const step = (y / 53) * 0.1;
      setZoom(zoom - step < 0.25 ? 0.25 : zoom - step);
    },
    [zoom, setZoom]
  );
  const wheelBind = useWheel(wheelHandler, { domTarget: window });
  const handleZoomOut = React.useCallback(
    () => setZoom(zoom - 0.1 < 0.25 ? 0.25 : zoom - 0.1),
    [zoom, setZoom]
  );
  const handleZoomIn = React.useCallback(() => setZoom(zoom + 0.1), [
    zoom,
    setZoom,
  ]);

  const [editedSignatory, setEditedSignatory] = React.useState(null);
  const [editedSignatoryCrop, setEditedSignatoryCrop] = React.useState({});
  const handleSetDefaultCrop = React.useCallback(() => {
    setEditedSignatoryCrop({
      unit: "%",
      x: 10,
      y: 10,
      width: 20,
      height: 20,
    });
  }, [setEditedSignatoryCrop]);
  const handleEditSignatoryCrop = React.useCallback(
    (_, crop) => {
      setEditedSignatoryCrop(crop);
    },
    [setEditedSignatoryCrop]
  );

  const dragHandler = React.useCallback(
    ({ delta: [x, y] }) => {
      if (editedSignatory !== null) {
        return;
      }
      setPosX(posX + x);
      setPosY(posY + y);
    },
    [posX, setPosX, posY, setPosY, editedSignatory]
  );
  const dragBind = useDrag(dragHandler, { domTarget: window });

  const [signatories, setSignatories] = React.useState([]);
  const handleCreateSignatory = React.useCallback(() => {
    const name = window.prompt("Please name the signatory");
    if (!name) {
      return;
    }
    signatories.push({ name, page });
    setSignatories([...signatories]);
    setEditedSignatory(signatories.length - 1);
    handleSetDefaultCrop();
  }, [signatories, page, setSignatories, handleSetDefaultCrop]);
  const handleEditSignatory = React.useCallback(
    (index) => {
      if (!signatories[index]) {
        return;
      }
      const name = window.prompt(
        `Please choose a new name for ${signatories[index].name}`
      );
      if (!name) {
        return;
      }
      signatories[index].name = name;
      signatories[index].page = page;
      setSignatories([...signatories]);
      setEditedSignatory(index);
      handleSetDefaultCrop();
    },
    [signatories, page, setSignatories, handleSetDefaultCrop]
  );
  const handleViewSignature = React.useCallback(
    (index) => {
      if (!signatories[index]) {
        return;
      }
      setPage(signatories[index].page);
    },
    [signatories, setPage]
  );
  const handleDeleteSignatory = React.useCallback(
    (index) => {
      if (!signatories[index]) {
        return;
      }
      signatories.splice(index, 1);
      setSignatories([...signatories]);
      if (editedSignatory === index) {
        setEditedSignatory(null);
      }
    },
    [signatories, setSignatories, editedSignatory, setEditedSignatory]
  );
  const pageRef = React.useRef();
  const [pageProxy, setPageProxy] = React.useState(null);
  const handleRenderSuccess = React.useCallback(
    (proxy) => {
      setPageProxy(proxy);
      if (!pageRef.current) {
        return;
      }
      pageRef.current.dispatchEvent(
        new Event("medialoaded", { bubbles: true })
      );
    },
    [pageRef]
  );
  const handleCancelEditedSignatory = React.useCallback(() => {
    handleDeleteSignatory(editedSignatory);
  }, [handleDeleteSignatory, editedSignatory]);
  const handleSaveEditedSignatory = React.useCallback(() => {
    signatories[editedSignatory].x = editedSignatoryCrop.x;
    signatories[editedSignatory].y = editedSignatoryCrop.y;
    signatories[editedSignatory].width = editedSignatoryCrop.width;
    signatories[editedSignatory].height = editedSignatoryCrop.height;
    setSignatories([...signatories]);
    setEditedSignatory(null);
  }, [
    editedSignatory,
    setEditedSignatory,
    signatories,
    setSignatories,
    editedSignatoryCrop,
  ]);
  const completedSignatories = React.useMemo(
    () => signatories.filter((s) => !isNaN(s.x)),
    [signatories]
  );
  const signatoriesOnThisPage = React.useMemo(
    () =>
      completedSignatories.filter((signatory) => {
        return signatory.page === page;
      }),
    [completedSignatories, page]
  );
  const handleSave = React.useCallback(() => {
    onSave(completedSignatories);
  }, [onSave, completedSignatories]);

  const signatoryPositions = React.useMemo(() => {
    return signatoriesOnThisPage.map((signatory) => ({
      fontSize: `${zoom}rem`,
      left: `${signatory.x}%`,
      top: `${signatory.y}%`,
      width: `${signatory.width}%`,
      height: `${signatory.height}%`,
    }));
  }, [signatoriesOnThisPage, zoom]);

  const documentComponent = (
    <Document file={document} onLoadSuccess={handleLoadSuccess}>
      {numPages > 0 && (
        <Page
          renderAnnotationLayer={false}
          renderTextLayer={false}
          pageNumber={page}
          height={size.height * zoom}
          onRenderSuccess={handleRenderSuccess}
          inputRef={(newRef) => (pageRef.current = newRef)}
        />
      )}
    </Document>
  );

  return (
    <div className="document-position" {...dragBind()}>
      <div
        className="document-document"
        {...wheelBind()}
        style={{
          left: posX,
          top: posY,
        }}
      >
        {editedSignatory === null ? (
          documentComponent
        ) : (
          <ReactCrop
            crop={editedSignatoryCrop}
            onChange={handleEditSignatoryCrop}
            renderComponent={documentComponent}
          />
        )}
        {!!pageRef.current &&
          signatoriesOnThisPage.map((signatory, i) => (
            <div
              className="document-placed"
              key={i}
              style={signatoryPositions[i]}
            >
              <span>{signatory.name}</span>
            </div>
          ))}
      </div>

      <div className="document-zoom">
        <button onClick={handleZoomOut} className="button">
          <i className="fa fa-search-minus" />
        </button>
        {Math.floor(zoom * 100)}%
        <button onClick={handleZoomIn} className="button">
          <i className="fa fa-search-plus" />
        </button>
      </div>
      <div className="document-page">
        <button onClick={handlePreviousPage} className="button">
          <i className="fa fa-chevron-left" />
        </button>
        {page} / {numPages}
        <button onClick={handleNextPage} className="button">
          <i className="fa fa-chevron-right" />
        </button>
      </div>
      <div className="document-back">
        <button onClick={handleBack}>Upload another PDF</button>
      </div>
      {editedSignatory === null && completedSignatories.length > 0 && (
        <div className="document-save">
          <button onClick={handleSave}>Save this document</button>
        </div>
      )}
      <div className="document-palette">
        <h3>
          <span>Signatories</span>
          <button onClick={handleCreateSignatory} className="button">
            <i className="fa fa-plus" />
          </button>
        </h3>
        {signatories.map((person, i) => (
          <div
            key={i}
            className={clsx("document-signatory", {
              "document-signatory-placed": person.x !== undefined,
            })}
          >
            <span>{person.name}</span>
            {editedSignatory === i ? (
              <>
                <button onClick={handleSaveEditedSignatory} className="button">
                  <i className="fa fa-check" />
                </button>
                <button
                  onClick={handleCancelEditedSignatory}
                  className="button"
                >
                  <i className="fa fa-trash" />
                </button>
              </>
            ) : (
              <>
                <button
                  onClick={() => handleViewSignature(i)}
                  className="button"
                >
                  <i className="fa fa-eye" />
                </button>
                <button
                  onClick={() => handleEditSignatory(i)}
                  className="button"
                >
                  <i className="fa fa-pen" />
                </button>
                <button
                  onClick={() => handleDeleteSignatory(i)}
                  className="button"
                >
                  <i className="fa fa-times" />
                </button>
              </>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Position;
