"use client";

import React, { useContext } from "react";
import { AuthContext } from "./auth";
import { AppBar, Button, IconButton, Toolbar, Typography } from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";

const AppTopBar = () => {
  //const { currentUser } = useContext(AuthContext);
  //const currentUserEmail = currentUser ? currentUser.email : "";

  return (
    <AppBar>
      <Toolbar>
        <IconButton
          size="large"
          edge="start"
          color="inherit"
          aria-label="menu"
          sx={{ mr: 2 }}
        >
          <MenuIcon />
        </IconButton>
        <Typography variant="h6" component="div">
          Regncon 2023
        </Typography>
        <Button color="inherit">Login</Button>
      </Toolbar>
    </AppBar>
  );
};

export default AppTopBar;
