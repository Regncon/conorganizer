"use client";

import React from "react";
import Typography from "@mui/material/Typography";
import { FormControl, FormLabel, RadioGroup, FormControlLabel, Radio, FormControlLabelProps, styled, useRadioGroup, Card, CardContent, Divider } from "@mui/material";
import { ConEvent } from "@/lib/types";
import parse from 'html-react-parser';
import EventHeader from "./eventHeader";


interface Props {
    conEvent: ConEvent;
    showSelect?: boolean;
}

const EventUi = ({ conEvent, showSelect }: Props) => {

    interface StyledFormControlLabelProps extends FormControlLabelProps {
        checked: boolean;
    }

    const StyledFormControlLabel = styled((props: StyledFormControlLabelProps) => <FormControlLabel {...props} />)(
        ({ theme, checked }) => ({
            '.MuiFormControlLabel-label': checked && {
                color: theme.palette.primary.main,
            },
        })
    );

    function MyFormControlLabel(props: FormControlLabelProps) {
        const radioGroup = useRadioGroup();

        let checked = false;

        if (radioGroup) {
            checked = radioGroup.value === props.value;
        }

        return <StyledFormControlLabel checked={checked} {...props} />;
    }

    return (
        <Card
        sx={{ maxWidth: '440px' }}>
        <EventHeader conEvent={conEvent} />

        <Divider />
        <Typography className='pl-4 pr-4'
            sx={{minHeight: '7rem' }}        
        >
            {parse(conEvent?.description || '')}
            </Typography>

        <Divider />
        <CardContent
        sx={showSelect ?  { display: 'block' } : { display: 'none' }} 
        >
            <FormControl className="p-4">
                <FormLabel id="demo-row-radio-buttons-group-label">Puljepåmelding</FormLabel>
                <RadioGroup
                    row
                    aria-labelledby="demo-row-radio-buttons-group-label"
                    name="row-radio-buttons-group"
                    defaultValue="NotInterested"
                >
                    <MyFormControlLabel
                        value="NotInterested"
                        control={<Radio size="small" />}
                        label="Ikke intresert"
                    />
                    <MyFormControlLabel value="IfIHaveTo" control={<Radio size="small" />} label="Hvis jeg må" />
                    <MyFormControlLabel value="IWantTo" control={<Radio size="small" />} label="Har lyst" />
                    <MyFormControlLabel
                        value="RealyWantTo"
                        control={<Radio size="small" />}
                        label="Har veldig lyst"
                    />
                </RadioGroup>
            </FormControl>
        </CardContent>
        <Divider />
    </Card>
    );
};

export default EventUi;
