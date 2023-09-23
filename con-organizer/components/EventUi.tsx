'use client';

import { useEffect, useState } from 'react';
import { Box } from '@mui/material';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Divider from '@mui/material/Divider';
import FormControl from '@mui/material/FormControl';
import FormControlLabel, { FormControlLabelProps } from '@mui/material/FormControlLabel';
import FormLabel from '@mui/material/FormLabel';
import Radio from '@mui/material/Radio';
import RadioGroup, { useRadioGroup } from '@mui/material/RadioGroup';
import { styled } from '@mui/material/styles';
import Typography from '@mui/material/Typography';
import parse from 'html-react-parser';
import { ConEvent } from '@/models/types';
import EventHeader from './EventHeader';

type Props = {
    conEvent: ConEvent | undefined;
};

const EventUi = ({ conEvent }: Props) => {
    interface StyledFormControlLabelProps extends FormControlLabelProps {
        checked: boolean;
    }

    const [description, setDescription] = useState('');
    useEffect(() => {
        if (conEvent) {
            const tmp: string = conEvent?.description;
            // tmp = tmp.replace(/\n/g, '</p><p>');
            setDescription(tmp);
        }
    }, [conEvent]);
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
    // throw new Error(
    //     'lorem Ipsum error in conAuthor authorization dialog box - invalid Lorem, ipsum dolor sit amet consectetur adipisicing elit. Ea quia in blanditiis mollitia exercitationem, asperiores nam quidem commodi nulla illum laborum, distinctio magnam debitis vitae rerum, maiores maxime sapiente! Quia! Lorem, ipsum dolor sit amet consectetur adipisicing elit. Ea quia in blanditiis mollitia exercitationem, asperiores nam quidem commodi nulla illum laborum, distinctio magnam debitis vitae rerum, maiores maxime sapiente! Quia!'
    // );
    return (
        <Card>
            <EventHeader conEvent={conEvent} />
            <Divider />
            <Box className="p-4" sx={{ minHeight: '7rem', display: 'grid', gap: '.5rem' }}>
                {parse(description || '')}
            </Box>

            <Divider />
            <CardContent>
                <FormControl>
                    <FormLabel id="demo-row-radio-buttons-group-label">
                        <Typography variant="h6">Påmelding</Typography>
                    </FormLabel>
                    <RadioGroup
                        row
                        aria-labelledby="demo-row-radio-buttons-group-label"
                        name="row-radio-buttons-group"
                        defaultValue="NotInterested"
                        sx={{
                            display: 'grid',
                            width: '100vw',
                            maxWidth: '1080px',
                            padding: '.2em',
                            gridAutoFlow: 'column',
                            gridAutoColumns: '1fr',
                            placeContent: 'center',
                        }}
                    >
                        <MyFormControlLabel
                            sx={{ display: 'grid', textAlign: 'center', p: '.4em' }}
                            value="NotInterested"
                            control={<Radio size="small" />}
                            label="Ikke interessert"
                        />
                        <MyFormControlLabel
                            value="IfIHaveTo"
                            sx={{ display: 'grid', backgroundColor: '#00000055', textAlign: 'center', p: '.4em' }}
                            control={<Radio size="small" />}
                            label="Hvis jeg må"
                        />
                        <MyFormControlLabel
                            value="IWantTo"
                            sx={{ display: 'grid', backgroundColor: '#000000aa', textAlign: 'center', p: '.4em' }}
                            control={<Radio size="small" />}
                            label="Har lyst"
                        />
                        <MyFormControlLabel
                            value="RealyWantTo"
                            control={<Radio size="small" />}
                            label="Har veldig lyst"
                            sx={{ display: 'grid', backgroundColor: '#000000ff', textAlign: 'center', p: '.4em' }}
                        />
                    </RadioGroup>
                </FormControl>
            </CardContent>
            <Divider />
        </Card>
    );
};

export default EventUi;
