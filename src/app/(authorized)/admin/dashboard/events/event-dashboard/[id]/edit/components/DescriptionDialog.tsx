import {
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    Slide,
    TextField,
    Typography,
} from '@mui/material';
import { TransitionProps } from '@mui/material/transitions';
import { ReactElement, Ref, forwardRef, useEffect, useState } from 'react';
import MuiMarkdown from 'mui-markdown';
import { ConEvent } from '$lib/types';
import CancelIcon from '@mui/icons-material/Cancel';
import SaveIcon from '@mui/icons-material/Save';

type props = {
    data: ConEvent;
    handleSave: (data: ConEvent) => void;
    open: boolean;
    close: () => void;
};

const Transition = forwardRef(function Transition(
    props: TransitionProps & {
        children: ReactElement;
    },
    ref: Ref<unknown>
) {
    return <Slide direction="up" ref={ref} {...props} />;
});
const DescriptionDialog = ({ data, handleSave, close: Close, open }: props) => {
    const [description, setDescription] = useState(data.description);
    useEffect(() => {
        setDescription(data.description);
    }, [data, open]);
    return (
        <Dialog fullScreen open={open} TransitionComponent={Transition}>
            <Box
                sx={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fit,minmax(306px, 430px))',
                    gap: '16px',
                    justifyContent: 'center',
                    alignItems: 'center',
                }}
            >
                <DialogContent>
                    <Typography variant="h2">Rediger</Typography>
                    <DialogContentText>
                        Beskrivelse av arrangementet. Du kan bruke markdown for å formatere teksten.
                    </DialogContentText>
                    <TextField
                        autoFocus
                        fullWidth
                        multiline
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                    />
                </DialogContent>
                <DialogContent>
                    <Typography variant="h2">Forhåndsvisning</Typography>
                    <hr />
                    <MuiMarkdown>{description}</MuiMarkdown>
                </DialogContent>
                <Box />
                <DialogActions>
                    <Button variant="contained" color="error" startIcon={<CancelIcon />} onClick={Close}>
                        Cancel
                    </Button>
                    <Button
                        variant="contained"
                        color="info"
                        startIcon={<SaveIcon />}
                        onClick={() => {
                            data.description = description;
                            handleSave(data);
                        }}
                    >
                        save
                    </Button>
                </DialogActions>
            </Box>
        </Dialog>
    );
};
export default DescriptionDialog;
