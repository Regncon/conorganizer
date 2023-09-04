"use client";

import React, { useContext } from "react";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import Typography from "@mui/material/Typography";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { FormControl, FormLabel, RadioGroup, FormControlLabel, Radio, FormControlLabelProps, styled, useRadioGroup, Box, Card, CardHeader, CardContent, CardMedia, Button } from "@mui/material";
import { Edit } from "@mui/icons-material";
import EditDialog from "./editDialog";
import { serverTimestamp, doc, updateDoc, CollectionReference, DocumentData } from "firebase/firestore";
import { ConEvent } from "@/lib/types";

interface Props {
  conEvent: ConEvent;
  colletionRef: CollectionReference<DocumentData, DocumentData>
}


const EventUi = (props: Props) => {
  const [open, setOpen] = React.useState(false);

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

  const handleClose = () => {
    setOpen(false);
  };


  return (
    <Box className="event-ui">
      <EditDialog conEvent={props.conEvent} colletionRef={props.colletionRef} open={open} handleClose={handleClose} />
      <Card>
        <CardHeader
          title={props.conEvent?.title}
          subheader="Rom 222, Søndag kl 12:00 til 16:00"
        />
        <CardMedia
          component="img"
          height="194"
          image="/placeholder.jpg"
          alt={props.conEvent?.title}
        />
        <CardContent>
          <Typography>
            {props.conEvent?.description}
          </Typography>
        </CardContent>
        <hr />
        <FormControl className="p-4">
          <FormLabel id="demo-row-radio-buttons-group-label">Puljepåmelding</FormLabel>
          <RadioGroup
            row
            aria-labelledby="demo-row-radio-buttons-group-label"
            name="row-radio-buttons-group"
            defaultValue="NotInterested"
          >
            <MyFormControlLabel value="NotInterested" control={<Radio size="small" />} label="Ikke intresert" />
            <MyFormControlLabel value="IfIHaveTo" control={<Radio size="small" />} label="Hvis jeg må" />
            <MyFormControlLabel value="IWantTo" control={<Radio size="small" />} label="Har lyst" />
            <MyFormControlLabel value="RealyWantTo" control={<Radio size="small" />} label="Har veldig lyst" />
          </RadioGroup>
        </FormControl>
        <Button onClick={() => {
          setOpen(true);
        }}>Endre</Button>
      </Card>
    </Box>
  );
};

export default EventUi;
