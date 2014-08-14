package main
import "fmt"
import "strings"
import "errors"

func processGok(code string) (string, string, error) {
    var gofile string = "%s\nfunc %s(gok *Gok){\n%s\n}";
    funcName := fmt.Sprintf("Render%s", genRandName());
    imports, r, err := buildImports(code);
    if err != nil {
        return "", "", err;
    }
    goCode, err := buildGoCode(code[r:]);
    if err != nil {
        return "", "", err;
    }
    final := fmt.Sprintf(gofile, imports , funcName, goCode);
    return final, funcName, nil;
}

func buildImports(code string) (string, int, error) {
    p := "<?go-imports"
    pe := "?>";
    plen := len(p);
    indx := strings.Index(code, p);
    if indx == -1 {
        return "", 0, nil;
    }
    indxEnd := strings.Index(code, pe);
    if indxEnd == -1 {
        return "", -1, errors.New("unknown code pattern");
    }
    if indxEnd == (indx + plen) {
        return "", -1, nil;
    }
    imports := strings.Split(code[(indx+plen):indxEnd], "\n");
    for i := range imports {
        imports[i] = strings.TrimSpace(imports[i]);
    }
    return strings.Join(imports, "\n"), (indxEnd+2), nil;
}

func buildGoCode(code string) (string, error) {
    echoFunc := "\ngok.Echo(\"%s\");";
    p := "<?go ";
    pe := "?>";
    codeLen := len(code);
    result := make([]byte, 0, codeLen*3);
    for last := 0; last < codeLen; {
        slice := code[last:];
        indx := strings.Index(slice, p);
        if indx == -1 {
            if len(slice) > 0 {
                echo := fmt.Sprintf(echoFunc, strToCStr(slice));
                result = append(result, echo...);
            }
            break;
        }
        if indx != 0 {
            echo := fmt.Sprintf(echoFunc, strToCStr(slice[0:indx]))
            result = append(result, echo...);
        }
        indxEnd := strings.Index(slice, pe);
        if indxEnd == -1 {
            return "", errors.New("unknown code pattern");
        }
        if indxEnd == (indx+5) {
            last += indxEnd+2;
            continue;
        }
        goCode := strings.TrimSpace(slice[indx+5:indxEnd]);
        if len(goCode) != 0 {
            for _,v := range strings.Split(goCode, "\n") {
                result = append(result, '\n');
                result = append(result, strings.TrimSpace(v)...);
            }
        }
        last += indxEnd+2;
    }
    return string(result), nil;
}

func strToCStr(str string) string {
    result := strings.Replace(str, "\n", "\\n", -1);
    result = strings.Replace(result, "\t", "\\t", -1);
    result = strings.Replace(result, "\r", "\\r", -1);
    result = strings.Replace(result, "\v", "\\v", -1);
    result = strings.Replace(result, "\f", "\\f", -1);
    return result;
}