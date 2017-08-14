import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLMapper;
import org.json.JSONObject;
import org.json.XML;
import org.yaml.snakeyaml.Yaml;

import java.io.*;
import java.util.Map;

public class Convert {

    public static void main(String args[]){
        if(args==null){
            System.err.println("Can't Find args");
            return;
        }
        if(args[0] == null) {
            System.err.println("Can't Find args[0]");
            return;
        }
        if(args[1] == null) {
            System.err.println("Can't Find args[2]");
            return;
        }
        try {
            if (args[1].equals("json")) {
                xmltojson(args[0]);
            } else if (args[1].equals("yaml")) {
                xmltoyaml(args[0]);
            }
        }catch (Exception e){
            System.err.println(e.getMessage());
        }
    }

    public static void xmltojson(String xmlpath){
        JSONObject xmlJSONObj = XML.toJSONObject(readFile(xmlpath));
        String jsonPrettyPrintString = xmlJSONObj.toString(4);
        System.out.println(jsonPrettyPrintString);
    }

    public static void xmltoyaml(String xmlpath) throws  Exception{
        JSONObject xmlJSONObj = XML.toJSONObject(readFile(xmlpath));
        JsonNode jsonNodeTree = new ObjectMapper().readTree(xmlJSONObj.toString());
        String jsonAsYaml = new YAMLMapper().writeValueAsString(jsonNodeTree);
        Yaml yaml = new Yaml();
        Map<String, Object> data = (Map)yaml.load(jsonAsYaml);
        StringWriter writer;
        writer = new StringWriter();
        yaml.dump(data, writer);
        System.out.print(writer.toString());
    }

    static private String readFile(String FileName){
        File myFile=new File(FileName);
        if(!myFile.exists()) {
            System.err.println("Can't Find " + FileName);
        }
        try {
            BufferedReader in = new BufferedReader(new FileReader(myFile));
            StringBuilder builder = new StringBuilder();
            String str;
            while ((str = in.readLine()) != null) {
                builder.append(str);
            }
            in.close();
            return builder.toString();
        }
        catch (IOException e){
            e.getStackTrace();
        }
        return null;
    }
}
