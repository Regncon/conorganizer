import '@firebase/auth';
declare module '@firebase/auth' {
    interface ParsedToken {
        admin?: boolean;
    }
}
