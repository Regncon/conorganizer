"use client";

import React, { useState } from "react";
import {
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
  signOut,
} from "firebase/auth";
import { auth } from "../lib/firebase";
import Box from "@mui/material/Box";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import { Alert } from "@mui/material";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [success, setSuccess] = useState("");
  const [error, setError] = useState("");

  const register = () => {
    createUserWithEmailAndPassword(auth, email, password)
      .then(() => {
        resetInput();
      })
      .catch((err) => {
        console.error(err);
      });
  };

  const login = () => {
    signInWithEmailAndPassword(auth, email, password)
      .then(() => {
        resetInput();
        setSuccess("Logget inn! Ett øyeblikk...");
      })
      .catch((err) => {
        console.error(err);
        setError("Klarte ikke logge inn, ta kontakt hvis problemet vedvarer! Tekniske detaljer: " + err.message);
      });
  };

  const logOut = () => {
    signOut(auth);
  };

  const resetInput = () => {
    setEmail("");
    setPassword("");
  };
  return (
    <Box p={5} maxWidth={600} display={"grid"} justifyItems={"center"} gap={2}>
      <img src="/img/regnconlogony.png" alt="årets regncondrage" width={200} />
      <TextField
        label="e-post"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        fullWidth
      />
      <TextField
        label="passord"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        fullWidth
      />
      <Button variant="contained" mt={2} size="large" fullWidth onClick={login}>Logg inn</Button>
      { !!success && <Alert severity="success">{success}</Alert> }
      { !!error && <Alert severity="error">{error}</Alert> }

    </Box>
  );
};

export default Login;
