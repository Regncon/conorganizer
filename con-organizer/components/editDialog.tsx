"use client";

import { Button } from "../lib/mui";
import React, { useState } from "react";
import {
  doc,
  setDoc,
  serverTimestamp,
  CollectionReference,
  DocumentData,
  updateDoc,
} from "firebase/firestore";
import { Box, Dialog, DialogActions, DialogContent, DialogTitle, TextField } from "@mui/material";
import { ConEvent } from "@/lib/types";


interface Props {
  open: boolean;
  conEvent?: ConEvent;
  colletionRef: CollectionReference<DocumentData, DocumentData>;
  handleClose: () => void;
}

const EditDialog = ({ open, conEvent, colletionRef, handleClose }: Props) => {

  const [title, setTitle] = useState(conEvent?.title || "");
  const [description, setDescription] = useState(conEvent?.description || "");

  const addSchool = async () => {

    const newSchool = {
      title,
      description,
      createdAt: serverTimestamp(),
      lastUpdate: serverTimestamp(),
    };

    try {
      const schoolRef = doc(colletionRef);
      await setDoc(schoolRef, newSchool);
    } catch (error) {
      console.error(error);
    }
  };

  async function editEvent(conEvent: ConEvent) {
    const updatedSchool = {
      title: title,
      description: description,
      lastUpdate: serverTimestamp(),
    };

    try {
      const schoolRef = doc(colletionRef, conEvent.id);
      updateDoc(schoolRef, updatedSchool);
    } catch (error) {
      console.error(error);
    }
  }

  return (
    <Dialog
      open={open}
    >
      <Box
        sx={{ width: "900px", height: "900px" }}
      >
        <DialogTitle>{conEvent?.id ? "Endre" : "Legg til"}</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            id="title"
            label="Tittel"
            type="text"
            fullWidth
            variant="standard"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
          <TextField
            margin="dense"
            id="description"
            label="Beskrivelse"
            type="text"
            fullWidth
            multiline
            variant="standard"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          {conEvent?.id
            ? <Button onClick={() => editEvent(conEvent)}>Endre</Button>
            : <Button onClick={() => addSchool()}>Legg til</Button>
          }
          <Button onClick={handleClose}>Lukk</Button>
        </DialogActions>
      </Box>
    </Dialog>
  );
};

export default EditDialog;
