#CONVERT ALPHA-2 COUNTRIES TO ALPHA-3
#https://raw.githubusercontent.com/lukes/ISO-3166-Countries-with-Regional-Codes/master/all/all.json



import sys
import json

def main():
    if len(sys.argv) == 2:
        path = sys.argv[1]
    else:
        print("Need 1 argument: CSV IP Database")
        return ValueError

    try:
        countries_lookup = open("./countries_lookup.json", "r")
        json_countries = json.loads(countries_lookup.read())
    except FileNotFoundError:
        print("Couldn't find 'countries_lookup.json' or format was incorrect ")
                    #https://raw.githubusercontent.com/lukes/ISO-3166-Countries-with-Regional-Codes/master/all/all.json#
        return FileNotFoundError
    
    try:
        ip_database = open(sys.argv[1], "r+")
                    #https://db-ip.com/db/download/ip-to-country-lite
    except Exception as e:
        print("Couldn't open "+ sys.argv[1] + " : " + str(e))
        return e
    

    newLines = []
    for line in ip_database.readlines():
        # format:                       103.23.174.0,103.23.174.255,AU

        parts = line.split(",")
        # print(line, end="")
        if ":" in parts[1]: #entered ipv6 region. not usefull
            break
        if len(parts[2]) >= 4:
            print("File already in alpha-3 format")
            return 0
        for country in json_countries:
            if parts[2][:-1] == country["alpha-2"]:
                newLine = line[:-3] + country["alpha-3"]
                # print(newLine)
                newLines.append(newLine)
                break
    ip_database.close()

    try:
        ip_database = open(sys.argv[1], "w")
        for line in newLines:
            ip_database.write(f'{line}\n')
        ip_database.close()
    except Exception as e:
        print("Couldn't open "+ sys.argv[1] + " : " + str(e))
        return e
    
    return 0

if __name__ == "__main__":
    main()