"use client";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Button } from "../lib/mui";
import React, { useState } from "react";
import { faCoffee } from "@fortawesome/free-solid-svg-icons";

interface Props {}

const Experiment = (props: Props) => {
  const [clicked, setClicked] = useState(false);

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <Button
        variant="contained"
        onClick={() => setClicked(!clicked)}
        className="flex items-center gap-2"
      >
        hello World
      </Button>
      {clicked && <FontAwesomeIcon icon={faCoffee} />}
    </div>
  );
};

export default Experiment;
