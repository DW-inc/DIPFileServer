import SCSL.*;
import java.io.File;

public final class TestEnc
{
    public static void main(String[] args) {
        String srcFile,dstFile;

        srcFile = args[0];
        dstFile = "{ENC}" + args[0];

        SLDsFile sFile = new SLDsFile();
        SLBsUtil sUtil = new SLBsUtil();

            sFile.SettingPathForProperty("./softcamp.properties");
            sFile.SLDsInitDAC();

            sFile.SLDsAddUserDAC("SECURITYDOMAIN", "111001100", 0, 0, 0);

            int ret = sFile.SLDsEncFileDACV2("./keyDAC_SVR0.sc", "system", srcFile, dstFile, 1);
        
       }
}