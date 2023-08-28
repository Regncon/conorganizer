"use client";

import React, { useContext } from "react";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import Typography from "@mui/material/Typography";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";

const EventUi = () => {
  return (
    <Accordion className="">
      <AccordionSummary
        expandIcon={<ExpandMoreIcon />}
        aria-controls="panel1a-content"
        id="panel1a-header"
      >
        <Typography>Drager og Fangehull</Typography>
      </AccordionSummary>
      <AccordionDetails>
        <p>Rom 416</p>
        <Typography>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse
          malesuada lacus ex, sit amet blandit leo lobortis eget.
        </Typography>
        <p>Ikke intresert, hvis jeg må, har lyst, har super ultra mega</p>
      </AccordionDetails>
    </Accordion>
  );
};

export default EventUi;
