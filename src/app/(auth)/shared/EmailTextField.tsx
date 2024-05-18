import { TextField, InputAdornment } from '@mui/material';
import { emailRegExp } from './utils';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
type Props = {
    defaultValue?: string;
};

const EmailTextField = ({ defaultValue }: Props) => {
    return (
        <TextField
            type="email"
            name="email"
            autoComplete="email"
            label="e-post"
            variant="outlined"
            defaultValue={defaultValue}
            fullWidth
            required
            InputProps={{
                endAdornment: (
                    <InputAdornment position="end">
                        <AccountCircleIcon />
                    </InputAdornment>
                ),
            }}
            inputProps={{
                pattern: emailRegExp.source,
                title: 'epost@example.com',
            }}
        />
    );
};

export default EmailTextField;
