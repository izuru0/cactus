import React, { useEffect, useState } from "react";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import MenuItem from "@mui/material/MenuItem";
import Select from "@mui/material/Select";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import Alert from "@mui/material/Alert";
import {
  bridgeBackTokensBesu,
  getAssetReferencesBesu,
} from "../../api-calls/besu-api";

export default function BridgeBackDialog(props) {
  const [assetRefs, setAssetRefs] = useState([]);
  const [assetRefID, setAssetRefID] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [sending, setSending] = useState(false);

  useEffect(() => {
    async function fetchData() {
      const list = await getAssetReferencesBesu(props.user);
      setAssetRefs(list.filter((asset) => asset.recipient === props.user));
    }

    if (props.open) {
      setSending(false);
      setAssetRefID("");
      fetchData();
    }
  }, [props.open, props.user]);

  const handleChangeAssetRefID = (event) => {
    setAssetRefID(event.target.value);
  };

  const performBridgeBackTransaction = async () => {
    if (assetRefID === "") {
      setErrorMessage("Please choose a valid Asset Reference ID");
    } else {
      setSending(true);
      const amount = assetRefs.find(
        (asset) => asset.id === assetRefID,
      ).numberTokens;
      await bridgeBackTokensBesu(props.user, amount, assetRefID);
      props.onClose();
    }
  };

  return (
    <Dialog open={props.open} keepMounted onClose={props.onClose}>
      <DialogTitle>{"Bridge Back CBDC"}</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Select the {props.user}"s Asset Reference that represents the amount
          to bridge back.
        </DialogContentText>
        {assetRefs.length === 0 ? (
          <Alert severity="error">
            Must escrow tokens before trying to bridge back CBDC.
          </Alert>
        ) : (
          <Select
            fullWidth
            name="assetRefID"
            value={assetRefID}
            variant="outlined"
            defaultValue={assetRefID}
            onChange={handleChangeAssetRefID}
          >
            {assetRefs.map((asset) => (
              <MenuItem key={asset.id} value={asset.id}>
                {asset.id}
              </MenuItem>
            ))}
          </Select>
        )}
        {errorMessage !== "" && <Alert severity="error">{errorMessage}</Alert>}
      </DialogContent>
      <DialogActions>
        {sending ? (
          <Button disabled>Sending...</Button>
        ) : (
          <div>
            <Button onClick={props.onClose}>Cancel</Button>
            <Button onClick={performBridgeBackTransaction}>Confirm</Button>
          </div>
        )}
      </DialogActions>
    </Dialog>
  );
}
