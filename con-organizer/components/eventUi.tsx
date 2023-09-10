"use client";

import React from "react";
import Typography from "@mui/material/Typography";
import { FormControl, FormLabel, RadioGroup, FormControlLabel, Radio, FormControlLabelProps, styled, useRadioGroup, Box, Card, CardHeader, CardContent, CardMedia, Button, CardActions, Chip, Collapse, IconButton } from "@mui/material";
import { CollectionReference, DocumentData } from "firebase/firestore";
import { ConEvent } from "@/lib/types";
import parse from 'html-react-parser';
import { faDiceD20 } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import CloseIcon from '@mui/icons-material/Close';

interface Props {
  conEvent: ConEvent;
  colletionRef: CollectionReference<DocumentData, DocumentData>
}


const EventUi = ({ conEvent, colletionRef }: Props) => {
  const [open, setOpen] = React.useState(false);
  const [expanded, setExpanded] = React.useState(false);

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
    <Card>
      <Box onClick={() => {
        setExpanded(true);
      }} >

        <CardHeader
          sx="padding-bottom: 0.5rem;"
          title={conEvent?.title}
          subheader="Kjempebra spennende event."
        />

        <Box className="flex justify-start pb-4">
          <CardMedia
            className="ml-4"
            sx="width: 40%; max-height: 130px"
            component="img"
            image="/placeholder.jpg"
            alt={conEvent?.title}
          />
          <Box
            className="flex flex-col pl-4 pr-4" >
            <span>

              <Chip icon={<FontAwesomeIcon icon={faDiceD20} />}
                label="Rollespill"
                size="small"
                variant="outlined"
              />
            </span>
            <span>DnD 5e </span>                <span>
              Rom 222,
            </span>
            <span>
              Søndag: 12:00 - 16:00
            </span>
          </Box>
        </Box>
      </Box>
      <Collapse
        in={expanded}
      >
        <hr />
        <IconButton aria-label="lukk" 
        sx="background-color: gray;"
        onClick={() => {setExpanded(false)}}
        >
        <CloseIcon />
      </IconButton>

        <Typography>
          {parse(conEvent?.description || "")}
        </Typography>

        <hr />
        <CardContent>
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
        </CardContent>
        <hr />
        <CardActions>
          <Button onClick={() => {
            setOpen(true);
          }}>Endre</Button>
        </CardActions>
      </Collapse>
    </Card>
  );
};

export default EventUi;
