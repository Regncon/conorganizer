'use client';
import { ConEvent } from '@/lib/types';
import { faDiceD20 } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
    Card,
    Box,
    CardHeader,
    CardMedia,
    Chip,
    Collapse,
    IconButton,
    Typography,
    CardContent,
    FormControl,
    FormLabel,
    RadioGroup,
    Radio,
    CardActions,
    Button,
    useRadioGroup,
    styled,
    FormControlLabelProps,
    FormControlLabel,
    Divider,
} from '@mui/material';
import { collection, onSnapshot } from 'firebase/firestore';
import { useState, useEffect } from 'react';
import db from '../../../lib/firebase';
import parse from 'html-react-parser';
import CloseIcon from '@mui/icons-material/Close';

type Props = { id: string };

const Event = ({ id }: Props) => {
    const colletionRef = collection(db, 'schools');
    const [conEvents, setconEvents] = useState([] as ConEvent[]);
    const [loading, setLoading] = useState(false);
    useEffect(() => {
        setLoading(true);
        const unsub = onSnapshot(colletionRef, (querySnapshot) => {
            const items = [] as ConEvent[];
            querySnapshot.forEach((doc) => {
                items.push(doc.data() as ConEvent);
                items[items.length - 1].id = doc.id;
            });
            setconEvents(items);
            setLoading(false);
        });
        return () => {
            unsub();
        };
    }, []);

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
    const conEvent = conEvents.find((conEvent) => conEvent.id === id);
    return (
        <Card>
            <Box>
                <CardHeader
                    sx={{ paddingBottom: '0.5rem' }}
                    title={conEvent?.title}
                    subheader="Kjempebra spennende event."
                />

                <Box className="flex justify-start pb-4">
                    <CardMedia
                        className="ml-4"
                        sx={{ width: '40%', maxHeight: '130px' }}
                        component="img"
                        image="/placeholder.jpg"
                        alt={conEvent?.title}
                    />
                    <Box className="flex flex-col pl-4 pr-4">
                        <span>
                            <Chip
                                icon={<FontAwesomeIcon icon={faDiceD20} />}
                                label="Rollespill"
                                size="small"
                                variant="outlined"
                            />
                        </span>
                        <span>DnD 5e </span> <span>Rom 222,</span>
                        <span>Søndag: 12:00 - 16:00</span>
                    </Box>
                </Box>
            </Box>

            <Divider />
            <Typography>{parse(conEvent?.description || '')}</Typography>

            <Divider />
            <CardContent>
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
            <CardActions>
                <Button>Endre</Button>
            </CardActions>
        </Card>
    );
};

export default Event;
