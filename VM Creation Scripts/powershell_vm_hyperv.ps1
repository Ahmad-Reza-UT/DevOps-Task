# Define the VM name
$VM = "Kali Linux"

# Define the switch name
$Switch = "vEthernet"

# Define the install media path
$InstallMedia = "D:\ISO\kali-linux-2022.1-amd64.iso"

# Define the VM path
$VMPath = "D:\VM\$VM"

# Define the VHD path
$VHD = "$VMPath\$VM.vhdx"

# Creating a new VM
New-VM -Name $VM -MemoryStartupBytes 8GB -Path $VMPath -NewVHDPath $VHD -NewVHDSizeBytes 35GB -Generation 2 -SwitchName $Switch

# Adding the virtual DVD drive and mounting the ISO file
Add-VMDvdDrive -VMName $VM -Path $InstallMedia

# Configuring the boot order to DVD, VHD, and Network
Set-VMFirmware -VMName $VM -BootOrder $(Get-VMDvdDrive -VMName $VM), $(Get-VMHardDiskDrive -VMName $VM), $(Get-VMNetworkAdapter -VMName $VM)