import * as React from "react";
import axios from "axios";

import Drop from "./1-drop";
import Position from "./2-position";
import Complete from "./3-complete";

import { API_URL } from "../constants";

import "./style.css";

const Document = ({ username, id, hash }) => {
  // drop / position / complete
  const [mode, setMode] = React.useState("drop");

  // drop vars
  const [document, setDocument] = React.useState("");
  const onDrop = React.useCallback(
    (data) => {
      setDocument(data);
      setMode("position");
    },
    [setDocument]
  );

  const handleReupload = React.useCallback(() => setMode("drop"), [setMode]);
  const handleSave = React.useCallback(
    async (sigs) => {
      try {
        await axios({
          method: "post",
          url: API_URL + "/documents",
          data: document,
          params: {
            u: username,
            id: id,
            sig: hash,
            sigs: JSON.stringify(sigs.map(sig => ({
              x: sig.x / 100,
              y: sig.y / 100,
              width: sig.width / 100,
              height: sig.height / 100,
              name: sig.name,
            }))),
          },
        });
        setMode("complete");
      } catch (e) {
        alert(e.response.data.error);
      }
    },
    [document]
  );

  return (
    <div className="document-wrapper">
      <div className="document-modal">
        <div className="document-header">
          {mode === "drop"
            ? "Upload a new document"
            : mode === "position"
            ? "Edit signature fields"
            : mode === "Complete"
            ? "Done!"
            : "Invalid mode"}
        </div>
        <div className="document-body">
          {mode === "drop" && <Drop onDrop={onDrop} />}
          {mode === "position" && (
            <Position
              document={document}
              onBack={handleReupload}
              onSave={handleSave}
            />
          )}
          {mode === "complete" && <Complete username={username} id={id} />}
        </div>
      </div>
    </div>
  );
};

export default Document;
