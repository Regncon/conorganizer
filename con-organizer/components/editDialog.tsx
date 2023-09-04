"use client";

import { Button } from "../lib/mui";
import React, { useContext, useState } from "react";
import {
  doc,
  setDoc,
  collection,
  serverTimestamp,
} from "firebase/firestore";
import db from "../lib/firebase";
import { AuthContext } from "./auth";
import { Dialog, DialogTitle } from "@mui/material";

interface Props {
    open: boolean;
    handleClose: () => void;
}

const EditDialog = (props: Props) => {

  const colletionRef = collection(db, "schools");

  const { currentUser } = useContext(AuthContext);

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [score, setScore] = useState("");

  // ADD FUNCTION
  const addSchool = async () => {
    const owner = currentUser ? currentUser.uid : "unknown";
    const ownerEmail = currentUser ? currentUser.email : "unknown";

    const newSchool = {
      title,
      description,
      score: +score,
      owner,
      ownerEmail,
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

  return (
      <Dialog
        onClose={() => {
          console.log("close");
        }}
        open={props.open}
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
          <h6>Score 0-10</h6>
          <input
            className="text-black"
            type="number"
            value={score}
            onChange={(e) => setScore(e.target.value)}
          />
          <h6>Description</h6>
          <textarea
            className="text-black"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </div>
        <Button onClick={() => addSchool()}>Submit</Button>
        <Button onClick={props.handleClose}>Cancel</Button>
      </Dialog>
  );
};

export default EditDialog;
