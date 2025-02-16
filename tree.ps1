function Show-Tree {
    param (
        [string]$Path = (Get-Location),
        [switch]$ShowFiles
    )

    function Get-Tree {
        param ($Path, $Prefix)

        $folders = Get-ChildItem -Path $Path -Directory
        foreach ($folder in $folders) {
            Write-Output "$Prefix+-- $($folder.Name)"
            Get-Tree -Path $folder.FullName -Prefix ("$Prefix|   ")
        }

        if ($ShowFiles) {
            $files = Get-ChildItem -Path $Path -File
            foreach ($file in $files) {
                Write-Output "$Prefix+-- $($file.Name)"
            }
        }
    }

    Write-Output "$Path"
    Get-Tree -Path $Path -Prefix ''
}

# Run the function
Show-Tree -ShowFiles
