import SCSL.*;
import java.io.File;

public final class TestDec
{
    public static void main(String[] args) {
        String srcFile,dstFile;

        srcFile = args[0];
        dstFile = "{DEC}" + args[0];

        SLDsFile sFile = new SLDsFile();

        sFile.SettingPathForProperty("./softcamp.properties");
        sFile.SLDsInitDAC();

        int ret = sFile.CreateDecryptFileDAC (
		"./keyDAC_SVR0.sc",
		"SECURITYDOMAIN",
		srcFile,
		dstFile);
    }
}