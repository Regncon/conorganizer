"use client";

import React, { useContext } from "react";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import Typography from "@mui/material/Typography";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { FormControl, FormLabel, RadioGroup, FormControlLabel, Radio, FormControlLabelProps, styled, useRadioGroup } from "@mui/material";

const EventUi = () => {

  interface StyledFormControlLabelProps extends FormControlLabelProps {
    checked: boolean;
  }
  
  const StyledFormControlLabel = styled((props: StyledFormControlLabelProps) => (
    <FormControlLabel {...props} />
  ))(({ theme, checked }) => ({
    '.MuiFormControlLabel-label': checked && {
      color: theme.palette.primary.main,
    },
  }));
  
  function MyFormControlLabel(props: FormControlLabelProps) {
    const radioGroup = useRadioGroup();
  
    let checked = false;
  
    if (radioGroup) {
      checked = radioGroup.value === props.value;
    }
  
    return <StyledFormControlLabel checked={checked} {...props} />;
  }
  


  return (
    <Accordion className="">
      <AccordionSummary
        expandIcon={<ExpandMoreIcon />}
        aria-controls="panel1a-content"
        id="panel1a-header"
      >
        <Typography>Fangehull og Drager</Typography>
      </AccordionSummary>
      <AccordionDetails>
        <p>Rom 416</p>
        <Typography>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse
          malesuada lacus ex, sit amet blandit leo lobortis eget.
        </Typography>
        <FormControl>
          <FormLabel id="demo-row-radio-buttons-group-label">Puljepåmelding</FormLabel>
          <RadioGroup
            row
            aria-labelledby="demo-row-radio-buttons-group-label"
            name="row-radio-buttons-group"
            defaultValue="NotInterested"
          >
            <MyFormControlLabel value="NotInterested" control={<Radio size="small" />} label="Ikke intresert" />
            <MyFormControlLabel value="IfIHaveTo" control={<Radio size="small"/>} label="Hvis jeg må" />
            <MyFormControlLabel value="IWantTo" control={<Radio size="small"/>} label="Har lyst" />
            <MyFormControlLabel value="RealyWantTo" control={<Radio size="small"/>} label="Har veldig lyst" />
          </RadioGroup>
        </FormControl>
      </AccordionDetails>
    </Accordion>
  );
};

export default EventUi;
