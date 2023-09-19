'use client';

import { useEffect, useState } from 'react';
import {
    Card,
    CardContent,
    Divider,
    FormControl,
    FormControlLabel,
    FormControlLabelProps,
    FormLabel,
    Radio,
    RadioGroup,
    styled,
    useRadioGroup,
} from '@mui/material';
import Typography from '@mui/material/Typography';
import parse from 'html-react-parser';
import { ConEvent } from '@/models/types';
import EventHeader from './eventHeader';

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
            let tmp: string = conEvent?.description;
            tmp = tmp.replace(/\n/g, '</p><p>');
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

    return (
        <Card>
            <EventHeader conEvent={conEvent} />

            <Divider />
            <Typography variant="body1" className="p-4" sx={{ minHeight: '7rem', display: 'grid', gap: '.5rem' }}>
                {parse(description || '')}
            </Typography>

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
