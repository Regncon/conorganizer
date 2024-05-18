'use client';
import { faChevronLeft } from '@fortawesome/free-solid-svg-icons/faChevronLeft';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import IconButton from '@mui/material/IconButton';
import { useRouter } from 'next/navigation';

const BackButton = () => {
    const router = useRouter();
    return (
        <IconButton
            onClick={() => {
                router.back();
            }}
        >
            <FontAwesomeIcon icon={faChevronLeft} fixedWidth />
        </IconButton>
    );
};

export default BackButton;
