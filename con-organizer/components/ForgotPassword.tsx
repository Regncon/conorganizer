"use client";

import React, { useState } from "react";
import { sendPasswordResetEmail } from "firebase/auth";
import { auth } from "../lib/firebase";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import { Alert } from "@mui/material";

const ForgotPassword = (props: any) => {
  const [email, setEmail] = useState("");
  const [success, setSuccess] = useState("");
  const [error, setError] = useState("");
  const {setChoice} = props;

  const resetPwd = () => {
    sendPasswordResetEmail(auth, email)
      .then(() => {
        resetInput();
        setSuccess("Suksess! Vi har sendt deg en lenke for å skrive inn et nytt passord. Sjekk eposten din!");
      })
      .catch((err) => {
        console.error(err);
        setError("Klarte ikke sende epost, ta kontakt hvis problemet vedvarer! Tekniske detaljer: " + err.message);
      });
  };

  const resetInput = () => {
    setEmail("");
  };
  return (
    <Box p={5} maxWidth={600} display={"grid"} justifyItems={"center"} gap={1}>
      <h1>Glemt/endre passord</h1>
      <TextField
        label="e-post"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        fullWidth
      />
      <Button variant="contained" size="large" fullWidth onClick={resetPwd}>Send!</Button>
      <Button variant="outlined" size="large" fullWidth onClick={()=>setChoice("")}>Avbryt</Button>
      { !!success && <Alert severity="success">{success}</Alert> }
      { !!error && <Alert severity="error">{error}</Alert> }

    </Box>
  );
};

export default ForgotPassword;
