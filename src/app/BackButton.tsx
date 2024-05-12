'use client';
import ArrowBackOutlinedIcon from '@mui/icons-material/ArrowBackOutlined';
import IconButton from '@mui/material/IconButton';

const BackButton = () => {
    return (
        <IconButton
            onClick={() => {
                history.back();
            }}
        >
            <ArrowBackOutlinedIcon />
        </IconButton>
    );
};

export default BackButton;
