

Set objShell = CreateObject("Shell.Application")
Set fso = CreateObject("Scripting.FileSystemObject")
 
' Define the source file path and the destination path.
sourcePath = "ccs-client.exe"   
destinationPath = "C:\Program Files\ccs\"      
tempFolder = fso.GetSpecialFolder(0) 
testFilePath = tempFolder & "\ccs_test.txt"

' Check if you have administrator privileges.
On Error Resume Next
Set testFile = fso.CreateTextFile(testFilePath, True)
On Error GoTo 0

If Not fso.FileExists(testFilePath) Then
    ' If you do not have administrator privileges, request administrator privileges
    strPath = WScript.ScriptFullName
    objShell.ShellExecute "wscript.exe", Chr(34) & strPath & Chr(34), "", "runas", 1
	set fso = Nothing
	set objShell = Nothing
	WScript.Quit 1
Else
    testFile.Close
    fso.DeleteFile(testFilePath)
End If


If Not fso.FolderExists(destinationPath) Then
    fso.CreateFolder(destinationPath)
End If

Const adTypeBinary = 1
Const adSaveCreateOverWrite = 2

url = "http://10.10.10.10:56789/tools/ccs-client.exe"

savePath = CreateObject("Scripting.FileSystemObject").GetAbsolutePathName(".") & "\ccs-client.exe"

Set http = CreateObject("MSXML2.XMLHTTP")
http.Open "GET", url, False
http.Send

If http.Status = 200 Then
    Set stream = CreateObject("ADODB.Stream")
    stream.Type = adTypeBinary
    stream.Open
    stream.Write http.ResponseBody
    stream.SaveToFile savePath, adSaveCreateOverWrite
    stream.Close
Else
    MsgBox "Download failed, HTTP status£º" & http.Status
	WScript.Quit 1
End If

Set http = Nothing
Set stream = Nothing

If fso.FileExists(sourcePath) Then
	On Error Resume Next  
    fso.MoveFile sourcePath, destinationPath
	On Error GoTo 0  
	
Else
    MsgBox "Source file not found, please check the path.", vbExclamation, "Error"
    WScript.Quit 1
End If


shortcutName = "ccs-client.lnk"   

Set shell = CreateObject("WScript.Shell")
startupFolder = shell.SpecialFolders("Startup")

Set shortcut = shell.CreateShortcut(startupFolder & "\" & shortcutName)
shortcut.TargetPath = destinationPath & "\" &sourcePath
shortcut.WorkingDirectory = fso.GetParentFolderName(destinationPath)
shortcut.Save


On Error Resume Next

shell.Run Chr(34) & destinationPath & sourcePath & Chr(34), 0, False

WScript.Echo "Done!"

Set shell = Nothing
set fso = Nothing
set objShell = Nothing
