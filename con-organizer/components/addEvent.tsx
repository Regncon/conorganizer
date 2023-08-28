"use client";

import React, { useContext } from "react";
import { AuthContext } from "./auth";
import { Dialog, DialogTitle, IconButton } from "@mui/material";
import AddCircleIcon from "@mui/icons-material/AddCircle";

const AddEvent = () => {
  const [open, setOpen] = React.useState(false);
  return (
    <>
      <Dialog
        onClose={() => {
          console.log("close");
        }}
        open={open}
      >
        <DialogTitle>Legg til nytt arangement</DialogTitle>
      </Dialog>

      <IconButton
        className="absolute right-0 bottom-0"
        aria-label="add"
        color="error"
        size="large"
        onClick={() => {
          setOpen(true);
        }}
      >
        <AddCircleIcon />
      </IconButton>
    </>
  );
};

export default AddEvent;
