import SCSL.*;
import java.io.File;


public final class CheckEnc {
    public static void main (String[] args){
        String srcFile,dstFile;

        srcFile = args[0];
        dstFile = "{ENC}" + args[0];

        SLDsFile sFile = new SLDsFile();
        SLBsUtil sUtil = new SLBsUtil();
        
        int encrypted = sUtil.isEncryptFile(srcFile);
        System.out.println(encrypted);
    }
}
