"use client";

import React, { useContext } from "react";
import { AuthContext } from "./auth";

const Welcome = () => {
  const { currentUser } = useContext(AuthContext);
  const currentUserEmail = currentUser ? currentUser.email : "";

  return <h2>{`Welcome ${currentUserEmail}`}</h2>;
};

export default Welcome;
