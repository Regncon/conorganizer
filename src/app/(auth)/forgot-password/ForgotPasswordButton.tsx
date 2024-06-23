import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import { useFormStatus } from 'react-dom';

type Props = {};

const ForgotPasswordButton = ({}: Props) => {
    const { pending } = useFormStatus();
    return (
        <Button
            fullWidth
            type="submit"
            disabled={pending}
            endIcon={pending ? <CircularProgress size="1.5rem" /> : undefined}
        >
            Gløymd passord?
        </Button>
    );
};

export default ForgotPasswordButton;
