# GTC ChannelListMaker

A tool for creating an ordered and well-named channel list for a GTMedia GTC TV receiver.

### Preparing the tool

This tool uses Go-Lang to run. You might have to install it first so it will be required to compile the code.

Before compiling the code, you will need to have:
* A GTC exported channel list in XML format named **input.xml**
* A custom channel list named **desired.txt**. See "Generating the desired list" to know how to format the file.

### Generating the desired list

"desired.txt" has a very simple format.
* Every non-commented line will be representing a channel which number in the receiver will be the same as the line number. Comments are done by adding "//" to the beginning of the line.
* Every line's content will be the channel name. This name will be used to search inside the XML, so substantially modifying the name could confuse the search algorithm. In order to help the algorithm, it is possible to differentiate the display name and the search name using the following format:

```
Channel display name|Text to search the channel
```
For example, if we have a channel detected by the GTC which it's called "tdp" but we want to change its name to "Teledeporte", we would use:
```
Teledeporte|tdp
```
**An empty line will result in a fake channel called Sense nom (No Name)**, which is useful if you want to have some "blank spaces" between channels. In further versions I hope this "fake name" will be customizable.

For more examples, there is a file called "sample_desired.txt" with a complete and custom channel list.

### Compilation of the tool
The compilation of the code only requires the following commands.

```
git clone https://github.com/oscarsanchezdm/channellistmaker
cd channellistmaker
go build
```

### Running the tool
Once you have compiled the tool, just run it using
```
./channellistmaker
```
Please keep in mind that the tool requires to work the required files defined in "Preparing the tool".

### Menu items
The main screen is composed by a list of features
* **List XML channels**. This will show you the TV and radio channels that your exported XML file has. You can move along the list using the arrow keys. Pressing "Enter" key will give you audio and video properties.
* **Check desired channels**. This will generate your custom list. Every desired channel entry will show you the number of matches with the exported channel list, and pressing "Enter" key will show them in more detail.
* **Make a list** This is the option that exports the XML file to import it on the receiver. Firstly, it will check if there are channels with zero or more than a match, letting the user pick one. Then, the user will be asked for deleting the other channels, keeping the radio stations or keeping all of them.
* **Help**
* **Quit**

### Making the list
Before exporting the XML file, the tool will ensure that every channel of the designed list has only one match with the exported list. If there's more than a match for an entry, the tool will let the user pick one. A black screen will confirm that there are no duplicated matches on the list, so for exporting the channel list the user **will need to press the "Tab" key**. Then, the user will be asked for deleting the other channels, keeping the radio stations or keeping all of them.

### Exporting the list to your receiver
Once the tool has exported the XML file, copy it in a USB storage device and import it in your receiver going to Menu -> Installation -> DB Management -> Import. Then, you will have to go to Channel Configuration -> More -> Sort and select A-Z sorting. If the sortint does not have any effect, try using another type of sorting before selecting A-Z sorting again.

### Future features
These are some of the features that are expected to be added in a near future:
* Provide support to newer DB schemas
* Automatic desired channel list creation based in TV services channel lists

* Automatic import/export of channel list to the GTC receiver, using ADB
