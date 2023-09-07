"use client";

import { Button } from "../lib/mui";
import React, { useContext, useState } from "react";
import {
  doc,
  setDoc,
  collection,
  serverTimestamp,
  CollectionReference,
  DocumentData,
  updateDoc,
} from "firebase/firestore";
import db from "../lib/firebase";
import { AuthContext } from "./auth";
import { Dialog, DialogTitle } from "@mui/material";
import { ConEvent } from "@/lib/types";


interface Props {
  open: boolean;
  conEvent: ConEvent;
  colletionRef: CollectionReference<DocumentData, DocumentData>;
  handleClose: () => void;
}

const EditDialog = ( { open, conEvent, colletionRef, handleClose }: Props) => {  

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
      title : title,
      description: description,
      lastUpdate: serverTimestamp(),
    };

    console.log("editEvent", conEvent.id, updatedSchool);

    try {
      const schoolRef = doc(colletionRef, conEvent.id);
      updateDoc(schoolRef, updatedSchool);
    } catch (error) {
      console.error(error);
    }
  }

  return (
    <Dialog
      onClose={() => {
        handleClose();
      }}
      open={open}
    >
      <DialogTitle>Legg til nytt arangement</DialogTitle>

      <h1>Schools (SNAPSHOT adv.)</h1>
      <div className="inputBox">
        <h3>Add New</h3>
        <h6>Title</h6>
        <input
          className="text-black"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
        <h6>Description</h6>
        <textarea
          className="text-black"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
      </div>
      {conEvent?.id
        ? <Button onClick={() => editEvent(conEvent)}>Edit</Button>
        : <Button onClick={() => addSchool()}>Submit</Button>
      }
      <Button onClick={handleClose}>Cancel</Button>
    </Dialog>
  );
};

export default EditDialog;
